setTimeout(function () {
    location.reload();
}, 5000);

(function() {
    function number_format(num) {
        return num.toString().replace(/([0-9]+?)(?=(?:[0-9]{3})+$)/g , '$1,');
    }
    $('.number-format').text(function(i, e) {
        return number_format(e);
    })
})();
