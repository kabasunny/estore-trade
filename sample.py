import json
import os
import requests

# API のベースURL (本番環境、現行バージョン)
BASE_URL = "https://kabuka.e-shiten.jp/e_api_v4r5/"

# 認証情報（環境変数から取得することを推奨）
LOGIN_ID = os.environ.get("TACHIBANA_LOGIN_ID")
PASSWORD = os.environ.get("TACHIBANA_PASSWORD")


# ----------------------
# 1. ログイン処理
# ----------------------
def login():
    """ログインして必要なURLを取得する."""
    url = BASE_URL + "auth/login"
    payload = {
        "sCLMID": "CLMAuthLoginRequest",
        "sUserId": LOGIN_ID,
        "sPassword": PASSWORD,
    }

    try:
        response = requests.post(url, json=payload)
        response.raise_for_status()  # HTTPエラーをチェック
        data = response.json()

        if data["sResultCode"] == "0":
            print("ログイン成功")
            return {
                "request_url": data["sUrlRequest"],
                "master_url": data["sUrlMaster"],
                "price_url": data["sUrlPrice"],
                "event_url": data["sUrlEvent"],
                "session": requests.Session(),  # セッションを維持
            }
        else:
            print(f"ログイン失敗: {data['sResultCode']} - {data['sResultText']}")
            return None

    except requests.exceptions.RequestException as e:
        print(f"通信エラー: {e}")
        return None


# ----------------------
# 2. 株式新規注文
# ----------------------
def new_order(
    session,
    request_url,
    zyoutoeki_kazei_c,
    issue_code,
    baibai_kubun,
    order_price,
    order_suryou,
    second_password,
):
    """株式新規注文を行う."""
    url = request_url + "kabu/neworder"
    payload = {
        "sCLMID": "CLMKabuNewOrder",
        "sZyoutoekiKazeiC": zyoutoeki_kazei_c,  # 譲渡益課税区分 (1:特定, 3:一般, 5:NISA, 6:N成長)
        "sIssueCode": issue_code,  # 銘柄コード
        "sSizyouC": "00",  # 市場 (00:東証)
        "sBaibaiKubun": baibai_kubun,  # 売買区分 (1:売, 3:買)
        "sCondition": "0",  # 執行条件 (0:指定なし, 2:寄付, 4:引け, 6:不成)
        "sOrderPrice": order_price,  # 注文値段 (*:指定なし, 0:成行, その他:指値)
        "sOrderSuryou": order_suryou,  # 注文株数
        "sGenkinShinyouKubun": "0",  # 現金信用区分 (0:現物)
        "sOrderExpireDay": "0",  # 注文期日 (0:当日, その他:YYYYMMDD)
        "sGyakusasiOrderType": "0",  # 逆指値注文種別 (0:通常)
        "sGyakusasiZyouken": "0",  # 逆指値条件
        "sGyakusasiPrice": "*",  # 逆指値値段
        "sTatebiType": "*",  # 建日種類
        "sTategyokuZyoutoekiKazeiC": "*",  # 建玉譲渡益課税区分
        "sSecondPassword": second_password,  # 第二パスワード
    }

    try:
        response = session["session"].post(url, json=payload)
        response.raise_for_status()
        data = response.json()

        if data["sResultCode"] == "0":
            print(f"新規注文成功: 注文番号 = {data['sOrderNumber']}")
            return data
        else:
            print(f"新規注文失敗: {data['sResultCode']} - {data['sResultText']}")
            return None

    except requests.exceptions.RequestException as e:
        print(f"通信エラー: {e}")
        return None


# ----------------------
# 3. 株式訂正注文
# ----------------------
def correct_order(
    session, request_url, order_number, eigyou_day, order_price, second_password
):
    """株式訂正注文を行う."""
    url = request_url + "kabu/correctorder"
    payload = {
        "sCLMID": "CLMKabuCorrectOrder",
        "sOrderNumber": order_number,  # 注文番号 (訂正対象)
        "sEigyouDay": eigyou_day,  # 営業日 (訂正対象)
        "sCondition": "*",  # 執行条件 (*:変更なし, 0:指定なし, 2:寄付, 4:引け, 6:不成)
        "sOrderPrice": order_price,  # 注文値段 (*:変更なし, 0:成行, その他:指値)
        "sOrderSuryou": "*",  # 注文数量 (*:変更なし)
        "sOrderExpireDay": "*",  # 注文期日 (*:変更なし, 0:当日, その他:YYYYMMDD)
        "sGyakusasiZyouken": "*",  # 逆指値条件 (*:変更なし)
        "sGyakusasiPrice": "*",  # 逆指値値段 (*:変更なし)
        "sSecondPassword": second_password,  # 第二パスワード
    }

    try:
        response = session["session"].post(url, json=payload)
        response.raise_for_status()
        data = response.json()

        if data["sResultCode"] == "0":
            print("訂正注文成功")
            return data
        else:
            print(f"訂正注文失敗: {data['sResultCode']} - {data['sResultText']}")
            return None

    except requests.exceptions.RequestException as e:
        print(f"通信エラー: {e}")
        return None


# ----------------------
# 4. 株式取消注文
# ----------------------
def cancel_order(session, request_url, order_number, eigyou_day, second_password):
    """株式取消注文を行う."""
    url = request_url + "kabu/cancelorder"
    payload = {
        "sCLMID": "CLMKabuCancelOrder",
        "sOrderNumber": order_number,  # 注文番号 (取消対象)
        "sEigyouDay": eigyou_day,  # 営業日 (取消対象)
        "sSecondPassword": second_password,  # 第二パスワード
    }

    try:
        response = session["session"].post(url, json=payload)
        response.raise_for_status()
        data = response.json()

        if data["sResultCode"] == "0":
            print("取消注文成功")
            return data
        else:
            print(f"取消注文失敗: {data['sResultCode']} - {data['sResultText']}")
            return None

    except requests.exceptions.RequestException as e:
        print(f"通信エラー: {e}")
        return None


# ----------------------
# 5. メイン処理
# ----------------------
if __name__ == "__main__":
    # 1. ログイン
    login_info = login()
    if not login_info:
        exit()

    # 各APIのURL
    request_url = login_info["request_url"]

    # テスト用のパラメータ
    zyoutoeki_kazei_c = "1"  # 特定口座
    issue_code = "6758"  # ソニー
    baibai_kubun = "3"  # 買
    order_price = "0"  # 成行
    order_suryou = "100"  # 100株
    second_password = "your_second_password"  # 発注パスワード

    # 2. 新規注文
    new_order_result = new_order(
        login_info,
        request_url,
        zyoutoeki_kazei_c,
        issue_code,
        baibai_kubun,
        order_price,
        order_suryou,
        second_password,
    )

    if new_order_result:
        order_number = new_order_result["sOrderNumber"]
        eigyou_day = new_order_result["sEigyouDay"]

        # 3. 訂正注文 (指値に変更)
        correct_order_result = correct_order(
            login_info, request_url, order_number, eigyou_day, "12000", second_password
        )

        if correct_order_result:
            # 4. 取消注文
            cancel_order_result = cancel_order(
                login_info, request_url, order_number, eigyou_day, second_password
            )
