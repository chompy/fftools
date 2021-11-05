const FETCH_URL = "_data";
const JAIL_ELEMENT_ID = "jails"
const CHARACTER_ELEMENT_ID = "character";
let encounterId = 0;
let timeout = null;

function poll() {
    if (timeout) {
        clearTimeout(timeout);
    }
    fetch();
    timeout = setTimeout(poll, 2500);
}

function fetch() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4) {
            if (!this.responseText.trim()) {
                return;
            }
            var data = JSON.parse(this.responseText);
            setJails(data);
        }
    };
    xhttp.open("GET", FETCH_URL, true);
    xhttp.send();
}

function me() {
    return document.getElementById(CHARACTER_ELEMENT_ID).value;
}

function saveCharacter() {
    localStorage.setItem("titan_jails_character", me());
}
document.getElementById(CHARACTER_ELEMENT_ID).value = localStorage.getItem("titan_jails_character");

function tts(text) {
    let speaknow = new SpeechSynthesisUtterance(text); 
    window.speechSynthesis.speak(speaknow); 
}

function setJails(data) {
    if (data.encounter_id != encounterId && Object.keys(data.jails).length >= 3) {
        encounterId = data.encounter_id;
        jailElement = document.getElementById(JAIL_ELEMENT_ID);
        jailElement.innerHTML = "";
        let character = me();
        for (let k in data.jails) {
            let combatant = data.jails[k];
            let ele = document.createElement("li");
            ele.innerText = combatant;
            jailElement.append(ele);
            if (character == combatant) {
                tts(parseInt(k));
            }
        }
    }
}

poll();