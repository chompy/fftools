local decree_nisi_time = 3600
local is_tea_zone = false
local verdict_tracker = {}
local nisi_tracker = {}

local function on_zone(z)
    is_tea_zone = false
    if z == "The Epic Of Alexander (Ultimate)" then
        is_tea_zone = true
        log_info("Entered TEA.")
    end
end

local function on_encounter_change(e)
    nisi_tracker = {}
    verdict_tracker = {}
end

local function on_log(l)
    if not is_tea_zone then
        return
    end
    local match = regex_match("1A:([A-F0-9]*):(.*) gains the effect of Final (.*) (.) from|04:.*:Removing combatant Plasma Shield", l.log_line)
    if match ~= nil and match[1] ~= nil and match[1]:sub(0, 2) == "1A" then
        if match[4] == "Decree Nisi" then
            nisi_tracker[match[2]] = { match[5], act_encounter_time() }
        elseif match[4] == "Judgment: Decree Nisi" then
            verdict_tracker[match[2]] = match[5]
        end
    elseif match ~= nil and match[1] ~= nil and match[1]:sub(0, 2) == "04" and me().id then
        local my_id = int_to_hex(me().id)
        local my_nisi = nisi_tracker[my_id]
        local my_verdict = verdict_tracker[my_id]
        if my_nisi ~= nil then
            for k, v in pairs(verdict_tracker) do
                if k ~= my_id and v == my_nisi[1] and act_encounter_time() - nisi_tracker[k][2] > decree_nisi_time then
                    local partner = act_combatant_from_id(hex_to_int(k))
                    if partner ~= nil then
                        act_say(partner.name .. " knee see")
                    end
                    break
                end
            end
        else
            for k, v in pairs(nisi_tracker) do
                if k ~= my_id and v[1] == my_verdict and act_encounter_time() - v[2] < decree_nisi_time then
                    local partner = act_combatant_from_id(hex_to_int(k))
                    if partner ~= nil then
                        act_say(partner.name .. " knee see")
                    end
                    break
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