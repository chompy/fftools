local function on_log(l)
    if string.match(l.log_line, "Blazing path, lead me to iron rule!") then
        ffl_say("stack and out")
    elseif string.match(l.log_line, "Fleeting light! Amid a rain of stars, exalt you the red moon!") then
        ffl_say("spread and dive")
    elseif string.match(l.log_line, "Fleeting light! 'Neath the red moon, scorch you the earth!") then
        ffl_say("dive and stack")
    elseif string.match(l.log_line, "From on high I descend, the hallowed moon to call!") then
        ffl_say("dive and in")
    elseif string.match(l.log_line, "From on high I descend, the iron path to call!") then
        ffl_say("dive and out")
    elseif string.match(l.log_line, "From on high I descend, the iron path to walk!") then
        ffl_say("dive and out")
    elseif string.match(l.log_line, "O hallowed moon, shine you the iron path!") then
        ffl_say("in and out")
    elseif string.match(l.log_line, "O hallowed moon, take fire and scorch my foes!") then
        ffl_say("in and stack")
    elseif string.match(l.log_line, "Take fire, O hallowed moon!") then
        ffl_say("stack and in")
    elseif string.match(l.log_line, "From hallowed moon I descend, a rain of stars to bring!") then
        ffl_say("in, dive, and spread")
    elseif string.match(l.log_line, "From on high I descend, the moon and stars to bring!") then
        ffl_say("dive, in, spread")
    elseif string.match(l.log_line, "From hallowed moon I bare iron, in my descent to wield!") then
        ffl_say("in, out, dive")
    elseif string.match(l.log_line, "Unbending iron, descend with fiery edge!") then
        ffl_say("out, dive, stack")
    elseif string.match(l.log_line, "From hallowed moon I descend, upon burning earth to tread!") then
        ffl_say("in, dive, stack")
    elseif string.match(l.log_line, "Unbending iron, take fire and descend!") then
        ffl_say("out, stack, dive")
    end
end

function init()
    ffl_event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "UCOB Nael Quotes",
        desc = "Calls out the mechanics for the Nael quotes in UCOB."
    }
end