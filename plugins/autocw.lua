-- imagine you're tired of your friends talking about pokemon all the time,
-- well, here's the plugin for you!

local function add_cw(status)
   return (status.Text or ""):find("pokemon")
end

function plugin(status)
   if status.ContentWarning == "" and add_cw(status) then
      status.ContentWarning = "pokemon"
   end
end
