local decree_nisi_time = 36000
local is_tea_zone = false
local nisi_tracker = {}
local has_said = false

local function on_zone(z)
    is_tea_zone = false
    if z == "The Epic Of Alexander (Ultimate)" then
        is_tea_zone = true
        fft_log_info("Entered TEA.")
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
    fft_log_info("Local player needs to give/take Nisi from player #" .. tonumber(id, 16) .. ".")
    local partner = fft_combatant_from_id(tonumber(id, 16))
    if partner ~= nil then
        fft_say(partner.name .. " knee see")
    end
    has_said = true
end

local function player_nisi()
    for symbol, pid in pairs(nisi_tracker) do
        if pid == fft_me().id then
            return symbol
        end
    end
    return nil
end

local function on_log(l)
    local match = fft_regex_match("1A:([A-F0-9]*):(.*) gains the effect of Final (.*) (.) from", l.log_line)
    if match ~= nil then
        local match_pid = tonumber(match[2], 16)
        local match_symbol = match[5]
        if match[4] == "Decree Nisi" then
            nisi_tracker[match_symbol] = match_pid
        elseif match[4] == "Judgment: Decree Nisi" then
            local my_nisi = player_nisi()
            -- local player needs nisi
            if my_nisi == nil and match_pid == fft_me().id then
                say_partner(nisi_tracker[match_symbol])
            -- local player needs to give nisi
            elseif my_nisi ~= nil and my_nisi == match_symbol and match_pid ~= fft_me().id then
                say_partner(match_pid)
            end
        end
    end
end

function init()
    fft_event_attach("act:encounter:zone", on_zone)
    fft_event_attach("act:encounter:change", on_encounter_change)
    fft_event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "TEA Third Nisi",
        desc = "Calls the player's partner for third Nisi pass in The Epic of Alexander (Ultimate).",
        version = "0.01"
    }
end