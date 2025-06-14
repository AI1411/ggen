/**
 * Generated by orval v6.31.0 🍺
 * Do not edit manually.
 * 農業災害支援システム API
 * 農業災害の報告と支援申請を管理するためのAPI
 * OpenAPI spec version: 1.0
 */
import type { MyerrorsErrorCode } from './myerrorsErrorCode';
import type { MyerrorsErrorMessage } from './myerrorsErrorMessage';

export interface HandlerErrorResponse {
  /** 内部のエラーコード */
  code?: MyerrorsErrorCode;
  /** 内部のエラーメッセージ */
  message?: MyerrorsErrorMessage;
}
