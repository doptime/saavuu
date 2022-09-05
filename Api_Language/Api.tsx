import axios from "axios";
var msgpack = require('@ygoe/msgpack');
//convert key value list to form data
export const BuildFormData = (data: any) => {
    let formData = new FormData();
    for (let key in data) {
        formData.append(key, data[key]);
    }
    return formData;
}
const JwtRequest = (headers: any = {}) => {
    let jwt = localStorage.getItem("jwt");
    if (jwt == null) return axios.create({ headers });
    try {
        let jwto: any = JSON.parse(jwt);
        if (jwto !== null && !!jwto.jwt) headers["Authorization"] = jwto.jwt;
    } catch (e) {
    }
    return axios.create({ headers });
};
const SignOut = (e: any) => {
    var UnAuthorized = !!e.response && e.response.status === 401
    if (!UnAuthorized) return;
    var jwt = { jwt: "", sub: "", id: "", LastGetJwtTime: new Date().getTime() + Math.random() };
    localStorage.setItem("jwt", JSON.stringify(jwt))
}
export enum Action { GET, PUT, DELETE, }
export enum ResponseContentType { json = "application/json", jpeg = "image/jpeg", ogg = "audio/ogg", mpeg = "video/mpeg", mp4 = "video/mp4", none = "", text = "text/plain", stream = "application/octet-stream" }
export class ReqGet {
    constructor(public Key: string, public Field: string = "", public Queries: string = "", public RspType: ResponseContentType = ResponseContentType.json) {
        if (Queries === "") RspType = ResponseContentType.none
    }
}
export class ReqSet {
    constructor(public Service: string, public Queries: string = "", public RspType: ResponseContentType = ResponseContentType.json) {
        if (Queries === "") RspType = ResponseContentType.none
    }
}
const Url = "https://api.iam26.com:3080/rSvc"
export const RGet = (req: ReqGet, callback: Function = (d: any) => null) => {
    JwtRequest().get(`${Url}?Key=${req.Key}&Field=${req.Field}&Queries=${req.Queries}&RspType=${req.RspType}`)
        .then(rsb => callback(rsb.data)).catch(SignOut);
};
export const RDel = (req: ReqGet, callback: Function = (d: any) => null) => {
    JwtRequest().delete(`${Url}?Key=${req.Key}&Field=${req.Field}&Queries=${req.Queries}&RspType=${req.RspType}`)
        .then(rsb => callback(rsb.data)).catch(SignOut);
};
export const RSet = (req: ReqSet, data: object | FormData = {}, callback: Function = (d: any) => null) => {
    var header = Object.getPrototypeOf(data) === Object.prototype && Object.keys(data).length === 0 ? {} : { "Content-Type": "application/octet-stream" }
    JwtRequest(header).post(`${Url}?Service=${req.Service}&Queries=${req.Queries}&RspType=${req.RspType}`, msgpack.serialize(data))
        .then(rsb => callback(rsb.data)).catch(SignOut);
};