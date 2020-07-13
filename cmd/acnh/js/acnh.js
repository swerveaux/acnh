let acnh = function() {
    let checkboxes = document.querySelectorAll(".donated_checkbox");
    checkboxes.forEach(function(checkbox) {
        checkbox.addEventListener("click", function(self) {
            console.log("I got clicked!");
            console.log(self.currentTarget.dataset["name"]);
            bugStr = window.localStorage.getItem(self.currentTarget.dataset["critter_type"]);
            if (bugStr == null || bugStr == "") {
                buglist = [];
            } else {
                buglist = bugStr.split(",");
            }
            if (self.currentTarget.checked) {
                if (!buglist.includes(self.currentTarget.dataset["name"])) {
                    buglist.push(self.currentTarget.dataset["name"]);
                }
            } else {
                if (buglist.includes(self.currentTarget.dataset["name"])) {
                    buglist = removeItemAll(buglist, self.currentTarget.dataset["name"]);
                }
            }
            bugStr = buglist.join(",");
            window.localStorage.setItem(self.currentTarget.dataset["critter_type"], bugStr);
        })
    })

    // Start by unhiding all the bug rows.
    document.querySelectorAll(".bug_row").forEach(function(row) {
        row.classList.remove("hidden");
    })
};

function removeItemAll(arr, value) {
    var i = 0;
    while (i < arr.length) {
        if (arr[i] === value) {
            arr.splice(i, 1);
        } else {
            ++i;
        }
    }
    return arr;
}

acnh();
