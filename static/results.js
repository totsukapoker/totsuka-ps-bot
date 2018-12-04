'use strict';

const vm = new Vue({
    el: "#app",
    delimiters: ['[[', ']]'],
    data: {
        fetchedAt: '',
        isAvailable: false,
        game: {},
        results: [],
        transactions: [],
    },
    mounted: loadData(),
    updated: function() {
        $('.number-format').text(function(i, o) {
            return o.toString().replace(/([0-9]+?)(?=(?:[0-9]{3})+$)/g , '$1,');
        });
        $('.roi-format').each(function(i, o) {
            let addClass = 'text-dark'
            if ($(o).text() > 0) {
                addClass = 'text-primary';
            } else if ($(o).text() < 0) {
                addClass = 'text-danger';
            }
            $(o).addClass(addClass);
            // $(o).text($(o).text() + '%');
        });
    },
})

async function loadData() {
    let game = {};
    let transactions = [];
    let results = [];

    await axios.get('/api/games/1').then(function (response) {
        vm.fetchedAt = new Date().toLocaleString("ja");
        game = response.data;
        transactions = response.data.Transactions;
    })

    vm.game = game;

    if (vm.transactions.length < transactions.length) {
        transactions.forEach(function(transaction){
            let key = null;
            for (let i = 0; i < results.length; i++) {
                if (results[i]['User']['ID'] == transaction.UserID) {
                    key = i;
                    break;
                }
            }
            if (key != null) {
                results[key]['Amount']['All'] += transaction.Amount;
                if (transaction.IsBuyin == true) {
                    results[key]['Amount']['Buyin'] += transaction.Amount;
                }
            } else {
                let r = {
                    User: transaction.User,
                    Amount: {
                        All: transaction.Amount,
                        Buyin: 0,
                    },
                    ROI: 0,
                }
                r.User.Name = transaction.User.DisplayName;
                if (transaction.User.MyName !== "") {
                    r.User.Name = transaction.User.MyName;
                }
                if (transaction.IsBuyin == true) {
                    r.Amount.Buyin = transaction.Amount;
                }
                results.push(r)
            }
        });
        for (let i = 0; i < transactions.length; i++) {
            transactions[i]['User']['Name'] = transactions[i]['User']['DisplayName'];
            if (transactions[i]['User']['MyName'] !== '') {
                transactions[i]['User']['Name'] = transactions[i]['User']['MyName'];
            }
            transactions[i]['Type'] = '現在額';
            if (transactions[i]['IsBuyin'] == true) {
                transactions[i]['Type'] = 'バイイン';
            }
        }
        for (let i = 0; i < results.length; i++) {
            if (results[i]['Amount']['Buyin'] > 0) {
                results[i]['ROI'] = results[i]['Amount']['All'] / results[i]['Amount']['Buyin'] * 100 - 100;
            }
        }
        vm.results = results;
        vm.transactions = transactions;
        vm.isAvailable = true;
    }
    setTimeout(loadData, 5000);
}

