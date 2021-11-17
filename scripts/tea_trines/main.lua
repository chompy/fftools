local has_start = false
local trine_tracker = ""
local trine_markers = {
    a = "center",
    b = "one",
    c = "three",
    d = "four"
}

local function say(from, to)
    fft_say(trine_markers[from] .. " to " .. trine_markers[to])
end

local function calculate()
    if #trine_tracker < 2 then
        return
    elseif #trine_tracker == 2 and not has_start then
        if trine_tracker == "yb" or trine_tracker == "by" then
            fft_say(trine_markers.c)
        elseif trine_tracker == "br" or trine_tracker == "rb" then
            fft_say(trine_markers.d)
        elseif trine_tracker == "yr" then
            fft_say(trine_markers.b)
        elseif trine_tracker == "ry" then
            fft_say(trine_markers.a)
        end
        has_start = true
        return
    elseif #trine_tracker == 3 and has_start then
        if trine_tracker == "ybr" or trine_tracker == "yrb" then
            fft_say("to " .. trine_markers.d)
        elseif trine_tracker == "rby" or trine_tracker == "ryb" then
            fft_say("to " .. trine_markers.c)
        elseif trine_tracker == "byr" then
            fft_say("to " .. trine_markers.a)
        elseif trine_tracker == "bry" then
            fft_say("to " .. trine_markers.b)
        end
        trine_tracker = ""
        has_start = false
    end
end

local function on_log(l)
    if l.type ~= 0x15 then
        return
    end
    -- reset
    if l.ability_id == 0x488E then
        trine_tracker = ""
        has_start = false
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
    fft_event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "TEA Trine Callouts",
        desc = "Calls starting and ending spots for trines, using standard markers with C-1-3-4 in The Epic of Alexander (Ultimate).",
        version = "0.01"
    }
end