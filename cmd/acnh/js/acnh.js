var acnh = function() {
    var checkboxes = document.querySelectorAll(".donated_checkbox");
    checkboxes.forEach(function(checkbox) {
        checkbox.addEventListener("click", function(self) {
            console.log("I got clicked!");
        })
    })
}();