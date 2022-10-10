import axios from "axios";
var msgpack = require('@ygoe/msgpack');
const JwtRequest = (headers: any = {}) => {
    let jwt = localStorage.getItem("Authorization");
    if (!!jwt) headers["Authorization"] = jwt;
    return axios.create({ headers });
};
const SignOut = (e: any) => {
    var UnAuthorized = !!e.response && e.response.status === 401
    if (!UnAuthorized) return;
    let jwt = "", sub = "", id = "", LastGetJwtTime = new Date().getTime()
    localStorage.setItem("jwt", JSON.stringify({ jwt, sub, id, LastGetJwtTime }));
}
export enum Action { GET, PUT, DELETE, }
export enum RspType { json = "&RspType=application/json", jpeg = "&RspType=image/jpeg", ogg = "&RspType=audio/ogg", mpeg = "&RspType=video/mpeg", mp4 = "&RspType=video/mp4", none = "", text = "&RspType=text/plain", stream = "&RspType=application/octet-stream" }
const Url = "https://api.iam26.com:3080/rSvc"
export enum Cmd { HEXISTS = "HEXISTS", HGET = "HGET", HGETALL = "HGETALL", HMGET = "HMGET" }
export const GetUrl = (cmd = Cmd.HGET, Key: string, Field: string = "", rspType: RspType = RspType.json, Queries: string = "") =>
    `${Url}?Cmd=${cmd}&Key=${Key}&Field=${Field}&Queries=${Queries}${rspType}`
export const RGet = (cmd = Cmd.HGET, Key: string, Field: string = "", rspType: RspType = RspType.json, Queries: string = "", callback: Function = (d: any) => null) => {
    if (Queries === "") rspType = RspType.none
    JwtRequest().get(`${Url}?Cmd=${cmd}&Key=${Key}&Field=${Field}&Queries=${Queries}${rspType}`)
        .then(rsb => callback(rsb.data)).catch(SignOut);
};
export const RDel = async (Key: string, Field: string = "", rspType: RspType = RspType.json, Queries: string = "", callback: Function = (d: any) => null) => {
    JwtRequest().delete(`${Url}?Key=${Key}&Field=${Field}&Queries=${Queries}&RspType=${rspType}`)
        .then(rsb => callback(rsb.data)).catch(SignOut);
};
export const RSet = async (Service: string, data: any, rspType: RspType = RspType.json, Queries: string = "", callback: Function = (d: any) => null) => {
    if (Queries === "") rspType = RspType.none
    JwtRequest({ "Content-Type": "application/octet-stream" }).post(`${Url}?Service=${Service}&Queries=${Queries}${rspType}`, msgpack.serialize(data))
        .then(rsb => callback(rsb.data)).catch(SignOut);
};