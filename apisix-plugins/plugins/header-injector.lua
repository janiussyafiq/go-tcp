local core = require("apisix.core")
local plugin_name = "header-injector"

local schema = {
    type = "object",
    properties = {
        headers = {
            type = "object",
            minProperties = 1,
        }
    },
    required = { "headers" }
}

local _M = {
    version = 0.1,
    priority = 1500,
    name = plugin_name,
    schema = schema
}

function _M.access(conf, ctx)
    -- Inject each header from config
    for header_name, header_value in pairs(conf.headers) do
        core.request.set_header(ctx, header_name, header_value)
        core.log.warn("Injected header: ", header_name, " = ", header_value)
    end

    return
end

return _M
