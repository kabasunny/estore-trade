from tabulate import tabulate


def calculate_interest(principal, annual_rate, days_held):
    """日割り金利を計算"""
    daily_rate = annual_rate / 100 / 365
    return principal * daily_rate * days_held


def highlight_text(text, should_highlight):
    """テキストをハイライトする（ANSIエスケープコードを使用）"""
    if should_highlight:
        return f"\033[43m{text}\033[0m"  # 背景を黄色に
    else:
        return text


def display_costs(
    principal_list, annual_rate, commission_fees_list, days_list, column_widths
):
    """コストを表示（カラム幅を指定可能、手数料を超える金利をハイライト）"""

    print(f"年間金利 : {annual_rate} %")

    headers = ["約定金額", "手数料"] + [f"{days}日" for days in days_list]

    # 各行のデータを準備
    table = []
    for i, principal in enumerate(principal_list):
        commission_fee = commission_fees_list["立花証券e支店"][i]
        interests = [
            calculate_interest(principal, annual_rate, days) for days in days_list
        ]
        row = [
            f"{principal:,}".rjust(column_widths[0]),
            f"{commission_fee}".rjust(column_widths[1]),
        ] + [
            highlight_text(
                f"{interest:,.2f}".rjust(w), interest > commission_fee
            )  # ハイライト
            for interest, w in zip(interests, column_widths[2:])
        ]
        table.append(row)

    # ヘッダーのフォーマット（間隔調整）
    formatted_headers = [
        header.center(width) for header, width in zip(headers, column_widths)
    ]

    # 表示
    print("  ".join(formatted_headers))  # ヘッダーを中央揃え
    for row in table:
        print("  ".join(row))  # データを右揃えにする


# データ設定
principal_list = [
    100_000,
    200_000,
    500_000,
    1_000_000,
    1_500_000,
    3_000_000,
    6_000_000,
    10_000_000,
]
commission_fees_list = {"立花証券e支店": [77, 99, 187, 341, 407, 473, 814, 869]}
annual_rate = 1.94
days_list = [1, 2, 3, 7, 14, 30, 90, 180, 365]

# カラム幅設定
column_widths = [12, 8] + [14] * len(days_list)

# 表示
display_costs(
    principal_list, annual_rate, commission_fees_list, days_list, column_widths
)
