local function split(s, delimiter)
    result = {};
    for match in (s..delimiter):gmatch("(.-)"..delimiter) do
        table.insert(result, match);
    end
    return result;
end

local function party_sort(a, b)
    if a.id == ffl_me().id then
        return true
    end
    local job_order = ffl_config_get("sort_order")
    ak = -1
    bk = -1
    for k, v in ipairs(job_order) do
        if string.upper(v) == string.upper(a.job) then
            ak = k
        end
        if string.upper(v) == string.upper(b.job) then
            bk = k
        end
    end
    return ak < bk
end

local function get_sorted_party()
    combatants = ffl_combatants()
    table.sort(combatants, party_sort)
    return combatants
end

local function get_combatant_index(combatant)
    for k, v in ipairs(get_sorted_party()) do
        if v == combatant or v.name == combatant or v.id == combatant then
            return k
        end
    end
    return -1
end

local function do_keypress(index)
    if index >= 1 then
        local keys = split(ffl_config_get("key_map")[index], "-")
        ffl_key_press(keys[1], keys[2], keys[3], keys[4])
    end
end

local function on_mark(combatant)
    do_keypress(get_combatant_index(combatant))
end

local function on_clear()
    local keys = split(ffl_config_get("key_map")[9], "-")
    ffl_key_press(keys[1], keys[2], keys[3], keys[4])
end

function init()
    ffl_event_attach("am:mark", on_mark)
    ffl_event_attach("am:clear", on_clear)
end

function info()
    return {
        name = "Auto Markers",
        desc = "Provides event for auto marking player characters."
    }
end