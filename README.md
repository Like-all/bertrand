bertrand
========

Bertrand is a simple accounting program for those who won't learn ledger.

#Building & installation

    git clone https://github.com/Like-all/bertrand.git
    cd bertrand
    go get github.com/droundy/goopt
    go build
    cp bertrand /somewhere/in/the/$PATH/

Also, you need to create config.json in `~/.bertrand` directory:

    {
        "bertrandFile": "/path/to/accounting/file.csv"
    }

#Usage

Bertrand stores data in csv files, which contains four columns: **DATE**;**ACCOUNT**;**AMOUNT**;**COMMENT**(optional)

Invoking bertrand without arguments results into output all of accounts, subaccounts and it's current values from beginning of your accounting story.

Bertrand has two modes: checkout mode for posting your expenses and log mode.
In checkout mode you can move your expenses values from one account to another. Optionally, you can set a date of posting:

    bertrand --checkout --from salary.work --to cash.pocket --date 2014-01-01 --amount 40000.00 --comment "WOOT!"

or even shorter:

    bertrand -c -f salary.work -t cash.pocket -d 2014-01-01 -a 40000.00 -C "WOOT!"

It results into two postings in csv file:

    2014-01-01;salary.work;-40000.00;WOOT!
    2014-01-01;cash.pocket;40000.00;WOOT!

Now you can see what's going on in your pocket:

    $ bertrand | grep pocket
    cash.pocket: 40000.00
    $

TODO: Improve this README & fix error handling
