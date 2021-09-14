# 微信支付第三版


> 微信支付第三版

#### 使用空格定义代码块

    pays := wxpay.PayOrderWx{
        AppID: "-", //
        MchID: MchId, //商户号
        SerialNo: MchSerialNumber,//证书编号
        PrivateKey:PrivateKey, //私钥
        Certificate: Certificate,//公钥
        NotifyURL: NotifyURL,//通知地址
    }

    //支付类型
    PAYJSAPI  小程序｜公众号  
    PAYNAVITIVE 支付适用于PC网站、实体店单品或订单、媒体广告支付等场景   
    PAYAPP 商户通过在移动端应用APP中集成开放SDK调起微信支付模块来完成支付  
    PAYH5  H5支付主要用于触屏版的手机浏览器请求微信支付的场景，方便从外部浏览器唤起微信支付

    res,err := pays.PlaceOrder(wxpay.PAYH5,"13242543324433311",100,openid,"测试支付","pay")
    fmt.Println(res)
    fmt.Println(err)


>> 持续更新中。。。。。。。。。。....


