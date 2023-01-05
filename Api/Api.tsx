import axios from "axios";
var msgpack = require('@ygoe/msgpack');
const JwtRequest = (headers: any = {}) => {
    let jwt = localStorage.getItem("Authorization");
    if (!!jwt) headers["Authorization"] = jwt;
    let req = axios.create({ headers });
    req.interceptors.request.use(
        (config: any) => {
            if (config.method === "post" || config.method === "put") {
                //if type of data is Object ,convert to object
                if (typeof config.data === "object" && !(config.data instanceof Array)) config.data = Object.assign({}, config.data);
                config.data = msgpack.encode(config.data);
                config.headers["Content-Type"] = "application/octet-stream";
            }
            return config;
        },
        (error: any) => {
            return Promise.reject(error);
        }
    );
    req.interceptors.response.use(
        (response: any) => {
            if ("data" in response) return response.data;
            debugger
            return response
        },
        (error: any) => {
            SignOut(error);
            return Promise.reject(error);
        }
    );
    return req
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

export const HEXISTS = (Key: string, Field: string = "") =>
    JwtRequest().get(`${Url}?Cmd=HEXISTS&Key=${Key}&Field=${Field}&Queries=*&RspType=application/json`)

export const HSET = (Key: string, Field: string = "", data: any, rspType: RspType = RspType.json, Queries: string = "") =>
    JwtRequest().put(`${Url}?Cmd=HSET&Key=${Key}&Field=${Field}&Queries=${Queries}${rspType}`, data)

export const HGET = (Key: string, Field: string = "", Queries: string = "*", rspType: RspType = RspType.json) =>
    JwtRequest().get(`${Url}?Cmd=HGET&Key=${Key}&Field=${Field}${!!Queries ? ("&Queries=" + Queries) : ""}${rspType}`)
export const HGETALL = (Key: string, Queries: string = "*", rspType: RspType = RspType.json) =>
    JwtRequest().get(`${Url}?Cmd=HGETALL&Key=${Key}${!!Queries ? ("&Queries=" + Queries) : ""}${rspType}`)
export const HKEYS = (Key: string) =>
    JwtRequest().get(`${Url}?Cmd=HKEYS&Key=${Key}&Queries=*&RspType=application/json`)
export const HMGET = (Key: string, Fields: string[] = [], Queries: string = "*") =>
    JwtRequest().get(`${Url}?Cmd=HMGET&Key=${Key}&Field=${Fields.join(",")}${!!Queries ? ("&Queries=" + Queries) : ""}&RspType=application/json`)

export const ZRange = (Key: string, Start: number, Stop: number, WITHSCORES: boolean) =>
    JwtRequest().get(`${Url}?Cmd=ZRANGE&Key=${Key}&Start=${Start}&Stop=${Stop}&WITHSCORES=${WITHSCORES}&Queries=*&RspType=application/json`)
export const ZRank = (Key: string, Member: string) =>
    JwtRequest().get(`${Url}?Cmd=ZRANK&Key=${Key}&Member=${Member}&Queries=*&RspType=application/json`)
export const ZRANGEBYSCORE = (Key: string, Min: number, Max: number, WITHSCORES: boolean) =>
    JwtRequest().get(`${Url}?Cmd=ZRANGEBYSCORE&Key=${Key}&Min=${Min}&Max=${Max}&WITHSCORES=${WITHSCORES}&Queries=*&RspType=application/json`)

export const SISMEMBER = (Key: string, Member: string) => JwtRequest().get(`${Url}?Cmd=SISMEMBER&Key=${Key}&Member=${Member}`)
export const HDEL = async (Key: string, Field: string = "", rspType: RspType = RspType.json, Queries: string = "") =>
    JwtRequest().delete(`${Url}?Cmd=HDEL&Key=${Key}&Field=${Field}&Queries=${Queries}${rspType}`)
export const Service = async (Service: string, data: any, Queries: string = "*", rspType: RspType = RspType.json) =>
    JwtRequest().post(`${Url}?Service=${Service}&Queries=${Queries}${rspType}`, data)

