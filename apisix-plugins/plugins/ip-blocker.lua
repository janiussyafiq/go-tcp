local core = require("apisix.core")
local plugin_name = "ip_blocker"

local schema = {
    type = "object",
    properties = {
        blocked_ips = {
            type = "array",
            items = { type = "string " },
            minItems = 1,
        }
    },
    required = { "blocked_ips" }
}

local _M = {
    version = 0.1,
    priority = 2000,
    name = plugin_name,
    schema = schema
}

function _M.access(conf, ctx)
    -- Get client IP
    local client_ip = core.request.get_remote_client_ip(ctx)

    -- Check if IP is in blocked list
    for _, blocked_ip in ipairs(conf.blocked_ips) do
        if client_ip == blocked_ip then
            core.log.warn("Blocked IP: ", client_ip)
            return 403, { message = "Access denied" }
        end
    end

    -- IP is allowed, continue
    return
end

return _M
