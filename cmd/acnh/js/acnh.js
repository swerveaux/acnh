let acnh = function() {
    let checkboxes = document.querySelectorAll(".donated_checkbox");
    let critters = {};
    ["bugs", "fishes", "sea_creatures", "umbrellas"].forEach(function(listName) {
        critters[listName] = [];
        let str = window.localStorage.getItem(listName);
        if (str != null && str != "") {
            critters[listName] = str.split(",");
        }
    });

    checkboxes.forEach(function(checkbox) {
        checkbox.addEventListener("click", function(self) {
            let critterType = checkbox.dataset["critter_type"];
            if (self.currentTarget.checked) {
                if (!critters[critterType].includes(self.currentTarget.dataset["name"])) {
                    critters[critterType].push(self.currentTarget.dataset["name"]);
                }
            } else {
                if (critters[critterType].includes(self.currentTarget.dataset["name"])) {
                    critters[critterType] = removeItemAll(critters[critterType], self.currentTarget.dataset["name"]);
                }
            }
            let str = critters[critterType].join(",");
            window.localStorage.setItem(self.currentTarget.dataset["critter_type"], str);
        })
    });

    document.getElementById("show_donated_bugs").addEventListener("click", function(e) {
        setDonatedBugsVisibility(e.currentTarget.checked);
    });
    document.getElementById("show_donated_fish").addEventListener("click", function(e) {
        setDonatedFishVisibility(e.currentTarget.checked);
    });
    document.getElementById("show_donated_sea_creatures").addEventListener("click", function(e) {
        setDonatedSCVisibility(e.currentTarget.checked);
    });
    document.getElementById("show_aquired_umbrellas").addEventListener("click", function(e) {
        setDonatedUmbrellaVisibility(e.currentTarget.checked);
    });

    document.querySelectorAll(".bug_row").forEach(function(row) {
        if (!critters["bugs"].includes(row.dataset["name"])) {
            row.classList.remove("hidden");
        } else {
            row.classList.add("donated");
            row.children[0].querySelector('input').checked = true;
        }
    })
    document.querySelectorAll(".fish_row").forEach(function(row) {
        if (!critters["fishes"].includes(row.dataset["name"])) {
            row.classList.remove("hidden");
        } else {
            row.classList.add("donated");
            row.children[0].querySelector('input').checked = true;
        }
    })
    document.querySelectorAll(".sea_creature_row").forEach(function(row) {
        if (!critters["sea_creatures"].includes(row.dataset["name"])) {
            row.classList.remove("hidden");
        } else {
            row.classList.add("donated");
            row.children[0].querySelector('input').checked = true;
        }
    })
    document.querySelectorAll(".umbrella_row").forEach(function(row) {
        if (!critters["umbrellas"].includes(row.dataset["name"])) {
            row.classList.remove("hidden");
        } else {
            row.classList.add("donated");
            row.children[0].querySelector('input').checked = true;
        }
    })
};

function setDonatedBugsVisibility(visible) {
    document.querySelectorAll("#bug_table .donated").forEach(function(elem) {
        if (visible) {
            elem.classList.remove("hidden");
        } else {
            elem.classList.add("hidden");
        }
    });
}

function setDonatedFishVisibility(visible) {
    document.querySelectorAll("#fish_table .donated").forEach(function(elem) {
        if (visible) {
            elem.classList.remove("hidden");
        } else {
            elem.classList.add("hidden");
        }
    });
}

function setDonatedSCVisibility(visible) {
    document.querySelectorAll("#sea_creature_table .donated").forEach(function(elem) {
        if (visible) {
            elem.classList.remove("hidden");
        } else {
            elem.classList.add("hidden");
        }
    });
}

function setDonatedUmbrellaVisibility(visible) {
    document.querySelectorAll("#umbrella_table .donated").forEach(function(elem) {
        if (visible) {
            elem.classList.remove("hidden");
        } else {
            elem.classList.add("hidden");
        }
    });
}

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
