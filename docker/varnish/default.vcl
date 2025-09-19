vcl 4.1;

backend default {
    .host = "127.0.0.1";
    .port = "8080";
}

acl purge {
  "localhost";
  "127.0.0.0"/24;
}

sub vcl_recv {
    if (req.url ~ "^/health/ping" || req.url ~ "^/.*/webhook" || req.url ~ "^/metrics") {
        return (pass);
    }

    if (req.method == "PURGE") {
        if (!client.ip ~ purge) {
            return (synth(405, "Not allowed."));
        }
        return (purge);
    }

    if (req.method == "BAN") {
        if (!client.ip ~ purge) {
            return (synth(405, "Not allowed."));
        }
        if (req.http.X-Cache-Tag) {
            ban("obj.http.Cache-Tag ~ " + req.http.X-Cache-Tag);
        }
        return (synth(200, "Ban added."));
    }

    return (hash);
}

sub vcl_hash {
    hash_data(req.url);

    if (req.http.Accept-Encoding) {
        hash_data(req.http.Accept-Encoding);
    }
    if (req.http.Accept) {
        hash_data(req.http.Accept);
    }
}

sub vcl_deliver {
    unset resp.http.Via;
    unset resp.http.X-Varnish;
    unset resp.http.Server;

    if (obj.hits > 0) {
        set resp.http.X-Cache = "HIT";
    } else {
        set resp.http.X-Cache = "MISS";
    }

    if (obj.ttl < 0s && obj.ttl + obj.grace > 0s) {
        set resp.http.X-Cache = "HIT-GRACE";
    }
}

sub vcl_backend_response {
    if (beresp.status >= 200 && beresp.status < 400) {
      set beresp.ttl = 1h;
        set beresp.grace = 3h;
        set beresp.keep = 4h;

        unset beresp.http.Set-Cookie;

    } else {
        set beresp.ttl = 0s;
        set beresp.grace = 0s;
      set beresp.uncacheable = true;
    }

    return (deliver);
}
