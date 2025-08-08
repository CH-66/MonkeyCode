/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

import { ContentType, RequestParams } from "./httpClient";
import { DomainCodeSnippet, V1GetContextReq, WebResp } from "./types";

/**
 * @description 为IDE端提供代码片段上下文检索功能，使用API Key认证。支持单个查询和批量查询。
 *
 * @tags CodeSnippet
 * @name PostGetContext
 * @summary IDE端上下文检索
 * @request POST:/api/v1/ide/codesnippet/context
 * @secure
 * @response `200` `(WebResp & {
    data?: (DomainCodeSnippet)[],

})` OK
 */

export const postGetContext = (
  request: V1GetContextReq,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainCodeSnippet[];
    }
  >({
    path: `/api/v1/ide/codesnippet/context`,
    method: "POST",
    body: request,
    secure: true,
    type: ContentType.Json,
    format: "json",
    ...params,
  });
