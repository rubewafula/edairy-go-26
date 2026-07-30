package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	initializers "github.com/rubewafula/edairy-go-26/internal/initializers"
	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/models"
	"github.com/rubewafula/edairy-go-26/internal/services"
	"gopkg.in/yaml.v3"
)

type seedConfig struct {
	AccountTypes []struct {
		Name string `yaml:"name"`
	} `yaml:"account_types"`
	Categories []struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	} `yaml:"categories"`
	Accounts []struct {
		Code     string `yaml:"code"`
		Name     string `yaml:"name"`
		Category string `yaml:"category"`
	} `yaml:"accounts"`
	PostingRules []struct {
		Type        string `yaml:"type"`
		Debit       string `yaml:"debit"`
		Credit      string `yaml:"credit"`
		Description string `yaml:"description"`
	} `yaml:"posting_rules"`
}

func main() {
	initializers.LoadEnvVariables()
	db.ConnectToDatabase()
	if err := services.EnsureLedgerSchema(); err != nil {
		log.Fatalf("schema: %v", err)
	}

	configPath := os.Getenv("FINANCE_SEED_CONFIG")
	if configPath == "" {
		configPath = filepath.Join("deploy", "finance", "chart-of-accounts.yaml")
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	var cfg seedConfig
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		log.Fatalf("parse yaml: %v", err)
	}

	typeIDs := map[string]uint64{}
	for _, t := range cfg.AccountTypes {
		var row models.AccountType
		if err := db.DB.Where("name = ?", t.Name).First(&row).Error; err != nil {
			row = models.AccountType{Name: t.Name}
			if err := db.DB.Create(&row).Error; err != nil {
				log.Fatalf("create type %s: %v", t.Name, err)
			}
		}
		typeIDs[t.Name] = row.ID
		fmt.Printf("✔ type %s\n", t.Name)
	}

	categoryIDs := map[string]uint64{}
	for _, c := range cfg.Categories {
		typeID := typeIDs[c.Type]
		var row models.AccountCategory
		err := db.DB.Where("name = ? AND account_type_id = ?", c.Name, typeID).First(&row).Error
		if err != nil {
			row = models.AccountCategory{Name: c.Name, AccountTypeID: typeID}
			if err := db.DB.Create(&row).Error; err != nil {
				log.Fatalf("create category %s: %v", c.Name, err)
			}
		}
		categoryIDs[c.Name] = row.ID
		fmt.Printf("✔ category %s\n", c.Name)
	}

	accountIDs := map[string]uint64{}
	for _, a := range cfg.Accounts {
		catID := categoryIDs[a.Category]
		var row models.Account
		err := db.DB.Where("account_code = ?", a.Code).First(&row).Error
		if err != nil {
			row = models.Account{
				AccountCode:       a.Code,
				Name:              a.Name,
				AccountCategoryID: catID,
				IsPostable:        true,
				IsActive:          true,
			}
			if err := db.DB.Create(&row).Error; err != nil {
				log.Fatalf("create account %s: %v", a.Code, err)
			}
		}
		accountIDs[a.Code] = row.ID
		fmt.Printf("✔ account %s %s\n", a.Code, a.Name)
	}

	for _, r := range cfg.PostingRules {
		debitID := accountIDs[r.Debit]
		creditID := accountIDs[r.Credit]
		if debitID == 0 || creditID == 0 {
			log.Printf("⚠ skip rule %s — account not found (debit=%s credit=%s)", r.Type, r.Debit, r.Credit)
			continue
		}
		var row models.TransactionPostingRule
		err := db.DB.Where("transaction_type = ?", r.Type).First(&row).Error
		if err != nil {
			row = models.TransactionPostingRule{
				TransactionType: r.Type,
				DebitAccountID:  debitID,
				CreditAccountID: creditID,
				Description:     r.Description,
			}
			if err := db.DB.Create(&row).Error; err != nil {
				log.Fatalf("create rule %s: %v", r.Type, err)
			}
		} else {
			db.DB.Model(&row).Updates(map[string]interface{}{
				"debit_account_id":  debitID,
				"credit_account_id": creditID,
				"description":       r.Description,
			})
		}
		fmt.Printf("✔ rule %s\n", r.Type)
	}

	fmt.Println("✅ Finance seed completed")
}
