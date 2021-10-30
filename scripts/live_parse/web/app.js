const IMG_URL_PREFIX = "img/";
const FETCH_URL = "_data";
const ROLE_TANK = "tank";
const ROLE_HEALER = "healer";
const ROLE_DPS = "dps";
var ROLE_TABLE = {
    [ROLE_HEALER]: ["sch", "whm", "ast"],
    [ROLE_TANK]: ["war", "drk", "gnb", "pld"]
};
var ENCOUNTER_ELEMENT = document.createElement("div");
var ENCOUNTER_ELEMENTS = {
    "zone" : document.createElement("span"),
    "state" : document.createElement("span"),
    "time" : document.createElement("span")
}
var PLAYER_ELEMENT = document.createElement("div");
var currentEncounter = "";
var currentEncTime = -1;
var currentEncActive = false;
var timeout = null;
var timerTimeout = null;

function init() {
    // setup encounter elment
    ENCOUNTER_ELEMENT.className = "encounter";
    document.getElementById("app").appendChild(ENCOUNTER_ELEMENT);
    // -- zone
    ENCOUNTER_ELEMENTS.zone.className = "zone"
    ENCOUNTER_ELEMENTS.zone.innerText = "-";
    ENCOUNTER_ELEMENT.appendChild(ENCOUNTER_ELEMENTS.zone);
    // -- state
    ENCOUNTER_ELEMENTS.state.className = "state";
    ENCOUNTER_ELEMENTS.state.innerText = "inactive";
    ENCOUNTER_ELEMENT.appendChild(ENCOUNTER_ELEMENTS.state);
    // -- time
    ENCOUNTER_ELEMENTS.time.className = "time";
    ENCOUNTER_ELEMENT.appendChild(ENCOUNTER_ELEMENTS.time);
    updateTimer();
    // player list div
    var playerEle = document.createElement("div");
    playerEle.id = "players";
    playerEle.className = "players";
    document.getElementById("app").appendChild(playerEle);
    // begin polling for data
    poll();
    tickTimer();
}

function poll() {
    if (timeout) {
        clearTimeout(timeout);
    }
    fetch();
    timeout = setTimeout(poll, 2500);
}

function reset() {
    document.getElementById("players").innerHTML = "";
    currentEncTime = -1;
}

function updateTimer() {
    if (currentEncTime < 0 || isNaN(currentEncTime)) {
        ENCOUNTER_ELEMENTS.time.innerText = "--:--";
        return;
    }
    var seconds = currentEncTime % 60;
    var minutes = Math.floor(currentEncTime / 60);
    ENCOUNTER_ELEMENTS.time.innerText = (minutes < 10 ? "0": "") + minutes + ":" + (seconds < 10 ? "0" : "") + seconds;
}

function tickTimer() {
    if (timerTimeout) {
        clearTimeout(timerTimeout);
    }
    if (currentEncTime >= 0 && currentEncActive) {
        currentEncTime++;
        updateTimer();
    }
    setTimeout(tickTimer, 1000);
}

function updateEncounter(data) {
    // reset
    if (data.encounter.start_time != currentEncounter) {
        reset();
        currentEncounter = data.encounter.start_time
    }
    // set active state
    ENCOUNTER_ELEMENT.classList.remove("active");
    ENCOUNTER_ELEMENTS.state.innerText = "inactive";
    currentEncActive = false;
    if (data.encounter.active) {
        ENCOUNTER_ELEMENT.classList.add("active");
        ENCOUNTER_ELEMENTS.state.innerText = "active";
        currentEncActive = true;
    }
    // update zone name
    ENCOUNTER_ELEMENTS.zone.innerText = typeof data.encounter.zone != "undefined" ? data.encounter.zone : "-";
    // update time
    var startTime = new Date(data.encounter.start_time * 1000)
    var endTime = new Date(data.encounter.end_time * 1000);
    if (currentEncActive) {
        endTime = new Date();
    }
    var diff = parseInt((endTime - startTime) / 1000);
    if (diff < 0) { 
        diff = 0;
    }
    if (currentEncTime < diff || !currentEncActive) {
        currentEncTime = diff;
        updateTimer();
    }
}

function updateCombatants(data) {
    let combatants = Object.values(data.combatants);
    combatants.sort(function(a, b) {
        return parseInt(a.damage) < parseInt(b.damage) ? 1 : -1;
    });
    for (var i = 0; i < combatants.length; i++) {
        updateCombatant(combatants[i], i+1);
    }
}

function buildColumn(rows) {
    var element = document.createElement("div");
    element.className = "col";
    for (var i in rows) {
        var rowElement = document.createElement("div");
        rowElement.className = "row";
        element.appendChild(rowElement)
        for (var j in rows[i]) {
            var valElement = document.createElement("div");
            valElement.className = "value " + rows[i][j];
            valElement.innerText = "-";
            rowElement.appendChild(valElement);
        }
    }
    return element;
}

function sanitizeNumeric(value) {
    if (isNaN(value) || !isFinite(value)) {
        return 0;
    }
    return value;
}

function updateCombatant(data, sort) {
    if (!data || !data.name) {
        return;
    }
    var nameId = data.name.toLowerCase().replace(" ", "-").replace("'", "");
    var element = document.getElementById("player-" + nameId);
    // create new element if not exists
    if (!element) {
        // create main element
        element = document.createElement("div");
        element.id = "player-" + nameId;
        element.className = "player job-" + data.job.toLowerCase() + " role-" + getCombatantRole(data) + " alive";
        document.getElementById("players").appendChild(element);
        // -- column 1
        var colOneEle = document.createElement("div");
        colOneEle.className = "value col job job-" + data.job.toLowerCase();
        element.appendChild(colOneEle);
        // --- job
        var jobEle = document.createElement("img")
        jobEle.className = "job-img";
        colOneEle.appendChild(jobEle);
        // -- column 2
        var colTwoEle = buildColumn([
            ["name"], ["damage"]
        ]);
        element.appendChild(colTwoEle);
        // -- column 3
        var colThreeEle = buildColumn([
            ["healing", "deaths"]
        ]);
        element.appendChild(colThreeEle);
        // -- column 4
        var colFourEle = buildColumn([
            ["critical-hits", "critical-heals"]
        ]);
        element.appendChild(colFourEle);
    }
    element.style.order = sort;
    element.style.webkitOrder = sort;

    // update job
    var jobEle = element.getElementsByClassName("job-img")[0];
    jobEle.title = data.job.toLowerCase();
    if (!jobEle.title) {
        jobEle.title = "lb";
    }
    jobEle.alt = jobEle.title;
    jobEle.src = IMG_URL_PREFIX + "jobs/" + jobEle.title + ".png";
    // update name
    var nameEle = element.getElementsByClassName("name")[0];
    nameEle.innerText = data.name;
    nameEle.title = nameEle.innerText;
    // update damage
    var damageEle = element.getElementsByClassName("damage")[0];
    damageEle.innerText = sanitizeNumeric(data.damage / currentEncTime).toFixed(2);
    damageEle.title = damageEle.innerText + " damage per second (" + data.damage + " total damage).";
    // update healing
    var healingEle = element.getElementsByClassName("healing")[0];
    healingEle.innerText = sanitizeNumeric(data.damage_healed / currentEncTime).toFixed(2);
    healingEle.title = healingEle.innerText + " healing per second (" + data.damage_healed + " total healing).";
    // update deaths
    var deathEle = element.getElementsByClassName("deaths")[0];
    deathEle.innerText = sanitizeNumeric(data.deaths);
    deathEle.title = sanitizeNumeric(data.deaths) + " deaths.";
    // update crit hits
    var critEle = element.getElementsByClassName("critical-hits")[0];
    var critPerc = sanitizeNumeric((data.critical_hits / data.hits) * 100).toFixed(1);
    if (data.critical_hits <= 0) {
        critPerc = "0.0";
    }
    critEle.innerText = critPerc + "%";
    critEle.title + critEle.innerText + " critical hits (" + data.critical_hits + " out of " + data.hits + ").";
    // update crit heals
    var critHeals = element.getElementsByClassName("critical-heals")[0];
    critPerc = sanitizeNumeric((data.critical_heals / data.heals) * 100).toFixed(1);
    if (data.critical_heals <= 0) {
        critPerc = "0.0";
    }
    critHeals.innerText = critPerc + "%";
    critHeals.title + critHeals.innerText + " critical heals (" + data.critical_heals + " out of " + data.heals + ").";
}

function getCombatantRole(data) {
    for (var i in ROLE_TABLE) {
        if (ROLE_TABLE[i].indexOf(data.job.toLowerCase()) != -1) {
            return i;
        }
    }
    return ROLE_DPS;
}

function parse(data) {
    data = JSON.parse(data);
    if (data.me.id > 0) {
        for (let id in data.combatants) {
            if (id == data.me.id) {
                data.combatants[id].name = data.me.name;
            }
        }
    }
    return data;
}

function fetch() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4) {
            if (!this.responseText.trim()) {
                return;
            }
            var data = parse(this.responseText);
            updateEncounter(data);
            updateCombatants(data);
        }
    };
    xhttp.open("GET", FETCH_URL, true);
    xhttp.send();
}

window.onload = init;