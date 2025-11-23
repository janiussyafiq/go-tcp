local core = require("apisix.core")
local plugin_name = "request-logger"

local _M = {
    version = 0.1,
    priority = 1000,
    name = plugin_name,
    schema = {}
}

function _M.access(conf, ctx)
    -- Get HTTP method and URI
    local method = core.request.get_method()
    local uri = ctx.var.uri

    -- Log them
    core.log.error("REQUEST: ", method, " ", uri)

    return
end

return _M
