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
        if match[4] == "Decree Nisi" then
            nisi_tracker[match[2]] = { match[5], act_encounter_time() }
        elseif match[4] == "Judgment: Decree Nisi" then
            local my_id = int_to_hex(me().id)
            local me_has_nisi = act_encounter_time() - nisi_tracker[my_id][2] > decree_nisi_time
            -- player got verdict and needs to take nisi
            if not me_has_nisi and match[2] == my_id then
                for pid, pdata in pairs(nisi_tracker) do
                    if pdata[1] == match[5] then
                        say_partner(pid)
                        return
                    end
                end
            -- player has nisi and needs to give it to another verdicted player
            elseif me_has_nisi then
                local my_nisi = nisi_tracker[my_id][1]
                if my_nisi == match[5] then
                    say_partner(match[2])
                    return
                end
            end
        end
    end
end

function init()
    event_attach("act:encounter:zone", on_zone)
    event_attach("act:encounter:change", on_encounter_change)
    event_attach("act:log_line", on_log)
end