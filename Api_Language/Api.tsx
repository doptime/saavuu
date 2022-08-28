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
    try {
        let jwt = localStorage.getItem("jwt");
        if (jwt !== null) {
            let jwto: any = JSON.parse(jwt);
            if (jwto !== null && jwto["jwt"] !== null) headers["Authorization"] = jwto["jwt"];
        }
    } catch (e) { }
    return axios.create({ headers });
};
const SignOut = (e: any) => {
    var jwt = { jwt: "", sub: "", id: "", LastGetJwtTime: new Date().getTime() + Math.random() };
    localStorage.setItem("jwt", JSON.stringify(jwt))
}
export enum Action { GET, PUT, DELETE, }
export enum ExpectResponse { json = "application/json", jpeg = "image/jpeg", ogg = "audio/ogg", mpeg = "video/mpeg", mp4 = "video/mp4" }
export class ReqGet {
    constructor(public Key: string, public Field: string = "", public Queries: string = "", public Expect: ExpectResponse = ExpectResponse.json) { }
}
export class ReqSet {
    constructor(public Service: string, public Queries: string = "", public Expect: ExpectResponse = ExpectResponse.json) { }
}
const Url = "/rSvc"
export const RGet = (req: ReqGet, setState: Function = (d: any) => null, dataTransform: Function = (d: any) => d) => {
    JwtRequest().get(`${Url}?Key=${req.Key}&Field=${req.Field}&Queries=${req.Queries}&Expect=${req.Expect}`)
        .then(rsb => setState(dataTransform(rsb.data))).catch(SignOut);
};
export const RDel = (req: ReqGet, setState: Function = (d: any) => null, dataTransform: Function = (d: any) => d) => {
    JwtRequest().delete(`${Url}?Key=${req.Key}&Field=${req.Field}&Queries=${req.Queries}&Expect=${req.Expect}`)
        .then(rsb => setState(dataTransform(rsb.data))).catch(SignOut);
};
export const RSet = (req: ReqSet, data: object | FormData = {}, setState: Function = (d: any) => null, dataTransform: Function = (d: any) => d) => {
    var header = Object.getPrototypeOf(data) === Object.prototype && Object.keys(data).length === 0 ? {} : { "Content-Type": "application/octet-stream" }
    JwtRequest(header).post(`${Url}?Service=${req.Service}&Queries=${req.Queries}&Expect=${req.Expect}`, msgpack.serialize(data))
        .then(rsb => setState(dataTransform(rsb.data))).catch(SignOut);
};