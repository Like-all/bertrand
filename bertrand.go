package main

import (
    "strconv"
    "time"
    "encoding/json"
    "strings"
    "sort"
    "os/user"
    "os"
    "bufio"
    "io/ioutil"
    "fmt"
    goopt "github.com/droundy/goopt"
)

var param_from = goopt.String([]string{"-f", "--from"}, "", "set \"from\" account")
var param_to = goopt.String([]string{"-t", "--to"}, "", "set \"to\" account")
var param_amount = goopt.String([]string{"-a", "--amount"}, "", "set amount")
var param_date = goopt.String([]string{"-d", "--date"}, "now", "set date")
var param_depth = goopt.Int([]string{"-D", "--depth"}, 5, "set from account")
var param_sincedate = goopt.String([]string{"-S", "--since-date"}, "start", "set SINCE date")
var param_todate = goopt.String([]string{"-T", "--to-date"}, "now", "set TO date")
var param_comment = goopt.String([]string{"-C", "--comment"}, "", "set comment")
var param_checkout_balance = goopt.Flag([]string{"-c", "--checkout"}, []string{"-b", "--balance"}, "checkout", "check balance")

type Config struct {
    BertrandFile string
}

func LoadConfig() (c *Config, err error) {
    var bfile []byte
    usr, _ := user.Current()
    cfgpath := usr.HomeDir + "/.bertrand/config.json"
    if bfile, err = ioutil.ReadFile(cfgpath); err != nil {
        return
    }
    c = new(Config)
    err = json.Unmarshal(bfile, c)
    return
}
func ReadTransactions (path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}
func WriteTransaction(transaction, path string) error {
    file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    _, err = file.WriteString(transaction)
    file.Close()
    return err
}
func GetAccounts (data []string, depth int) ([]string) {
    accounts := make([]string, len(data))
    for n, i := range data {
        fields := strings.Split(i, ";")
        longAccount := strings.Split(fields[1], ".")
        var strippedAccount []string
        strippedAccount = append(strippedAccount, longAccount[0])
        if depth > 1 {
            for j := 1; j < len(longAccount) && j < depth; j++ {
                strippedAccount = append(strippedAccount, longAccount[j])
            }
        }
        accounts[n] = strings.Join(strippedAccount, ".")
    }
    sort.Strings(accounts)
    var uniqAccounts []string
    uniqAccounts = append(uniqAccounts, accounts[0])
    for i := 1; i < len(accounts); i++ {
        if accounts[i] != accounts[i-1] {
            uniqAccounts = append(uniqAccounts, accounts[i])
        }
    }
    return uniqAccounts
}


func main() {
    goopt.Description = func() string {
        return "Bertrand - advanced accounting program."
    }
    goopt.Version = "0.1"
    goopt.Summary = "advanced accounting program"
    goopt.Parse(nil)

    const shortForm = "2006-01-02"
    cfgparams, _ := LoadConfig()
    megalines, _ := ReadTransactions(cfgparams.BertrandFile)
    accounts := GetAccounts(megalines, *param_depth)

    if *param_checkout_balance {
        if *param_date == "now" {
            *param_date = time.Now().Format(shortForm)
        }
        var newTransaction string
        newTransaction = *param_date + ";" + *param_from + ";" + "-" + *param_amount + ";" + *param_comment + "\n"
        if err := WriteTransaction(newTransaction, cfgparams.BertrandFile); err != nil {
            fmt.Println(err)
        }
        newTransaction = *param_date + ";" + *param_to + ";" + *param_amount + ";" + *param_comment + "\n"
        if err := WriteTransaction(newTransaction, cfgparams.BertrandFile); err != nil {
            fmt.Println(err)
        }
    } else {
        spent := make([]float64, len(accounts))
        for n, i := range accounts {
            spent[n] = 0.0
            for _, j := range megalines {
                fields := strings.Split(j, ";")
                t, _ := time.Parse(shortForm, fields[0])
                date := t.Unix()
                acc := fields[1]
                amount, _ := strconv.ParseFloat(fields[2], 64)
                var sinceDate, toDate int64
                if *param_sincedate == "start" {
                    t, _ = time.Parse(shortForm, *param_sincedate)
                    sinceDate = t.Unix()
                } else {
                    t, _ = time.Parse(shortForm, *param_sincedate)
                    sinceDate = t.Unix()
                }
                if *param_todate == "now" {
                    toDate = time.Now().Unix()
                } else {
                    t, _ = time.Parse(shortForm, *param_todate)
                    toDate = t.Unix()
                }
                if date >= sinceDate && date <= toDate && strings.Contains(acc, accounts[n]) {
                    spent[n] += amount
                }
            }
            fmt.Printf("%s: %.2f\n", i, spent[n])
        }
    }
}
