local trine_tracker = ""
local trine_markers = {
    a = "center",
    b = "one",
    c = "three",
    d = "four"
}

local function say(from, to)
    act_say(trine_markers[from] .. " to " .. trine_markers[to])
end

local function calculate()
    if #trine_tracker < 3 then
        return
    end
    log_info("Found trine '" .. string.upper(trine_tracker) .. ".'")
    if trine_tracker == "ybr" then
        say("c","d")
    elseif trine_tracker == "yrb" then
        say("b","d")
    elseif trine_tracker == "byr" then
        say("c","a")
    elseif trine_tracker == "bry" then
        say("d","b")
    elseif trine_tracker == "rby" then
        say("d","c")
    elseif trine_tracker == "ryb" then
        say("c","a")
    end
    trine_tracker = ""
end

local function on_log(l)
    if l.type ~= 0x15 then
        return
    end
    -- reset
    if l.ability_id == 0x488E then
        trine_tracker = ""
    -- add trine
    elseif l.ability_id == 0x488F and math.floor(l.source_x) == 100 then
        local ypos = math.floor(l.source_y)
        if ypos == 92 then
            trine_tracker = trine_tracker .. "y"
        elseif ypos == 100 then
            trine_tracker = trine_tracker .. "b"
        elseif ypos == 108 then
            trine_tracker = trine_tracker .. "r"
        end
    end
    -- calculate route
    calculate()
end

function init()
    event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "TEA Trine Callouts",
        desc = "Calls starting and ending spots for trines, using standard markers with C-1-3-4."
    }
end