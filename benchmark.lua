request = function()
    wrk.headers["X-API-KEY"] = "123456"
    wrk.method = "POST"
    param_value = math.random(1,1000000000)
    longURL = "https://zingnews.vn/da-tim-thay-co-gai-16-tuoi-mat-tich-khi-den-tphcm-post-" .. param_value
    wrk.headers["Content-Type"] = "application/json"
    wrk.body = '{"url": "'..longURL..'"}' ;
    return wrk.format("POST", path)
end