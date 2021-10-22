local decree_nisi_time = 36000
local is_tea_zone = false
local nisi_tracker = {}
local has_said = false

local function on_zone(z)
    is_tea_zone = false
    if z == "The Epic Of Alexander (Ultimate)" then
        is_tea_zone = true
        log_info("Entered TEA.")
    end
end

local function on_encounter_change(e)
    nisi_tracker = {}
    has_said = false
end

local function say_partner(id)
    if has_said then
        return
    end
    log_info("Local player needs to give/take Nisi from player #" .. hex_to_int(id) .. ".")
    local partner = act_combatant_from_id(hex_to_int(id))
    if partner ~= nil then
        act_say(partner.name .. " knee see")
    end
    has_said = true
end

local function on_log(l)
    local match = regex_match("1A:([A-F0-9]*):(.*) gains the effect of Final (.*) (.) from", l.log_line)
    if match ~= nil then
        local match_pid = match[2]
        local match_symbol = match[5]
        if match[4] == "Decree Nisi" then
            nisi_tracker[match_symbol] = match_pid
        elseif match[4] == "Judgment: Decree Nisi" then
            local my_id = int_to_hex(me().id)
            if match_pid == my_id and nisi_tracker[match_symbol] ~= my_id then
                say_partner(nisi_tracker[match_symbol])
            end
        end
    end
end

function init()
    event_attach("act:encounter:zone", on_zone)
    event_attach("act:encounter:change", on_encounter_change)
    event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "[TEA] Third Nisi",
        desc = "Calls the player's partner for third Nisi pass in The Epic Of Alexander (Ultimate)."
    }
end