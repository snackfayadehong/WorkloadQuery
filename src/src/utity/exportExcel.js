import * as XLSXStyle from "xlsx-js-style";
import _ from "lodash";
import { purchaseSummaryFileds, purchaseSummaryDetailFileds } from "../model/purchaseSummary";

const Title = {
    v: "采购订单",
    t: "s",
    s: { font: { bold: true, sz: 13, name: "宋体" }, alignment: { vertical: "center", horizontal: "center" }, fill: { fgColor: { rgb: "87CEEB" } } }
};

const TitleDetail = {
    v: "订单明细",
    t: "s",
    s: { font: { bold: true, sz: 13, name: "宋体" }, alignment: { vertical: "center", horizontal: "center" }, fill: { fgColor: { rgb: "ffff00" } } }
};

const ExportExcelUtity = (res, workBookName, type) => {
    let workBook = XLSXStyle.utils.book_new();

    if (type == "NodeliveredPurchase" && res.length != 0) {
        let SupplierNameGroups = _.groupBy(res, "SupplierName");

        for (let supplierName in SupplierNameGroups) {
            let supplierData = SupplierNameGroups[supplierName];
            let jsonWorksheet = XLSXStyle.utils.book_new();
            // 添加表头
            const header = {
                v: supplierName,
                t: "s",
                s: { font: { bold: true, sz: 18, name: "宋体" }, alignment: { vertical: "center", horizontal: "center" } }
            };

            XLSXStyle.utils.sheet_add_aoa(jsonWorksheet, [[header]]);
            const merge = [{ s: { r: 0, c: 0 }, e: { r: 0, c: 11 } }];
            jsonWorksheet["!merges"] = merge;
            // 遍历主记录和子记录
            supplierData.forEach(item => {
                // 主记录字段中文
                const translatedItem = {};
                // 处理主记录字段
                for (const key in item) {
                    if (purchaseSummaryFileds[key]) {
                        translatedItem[purchaseSummaryFileds[key]] = item[key];
                    }
                }
                XLSXStyle.utils.sheet_add_json(jsonWorksheet, [{ "": Title }], { origin: -1 }); // 小标题
                XLSXStyle.utils.sheet_add_json(jsonWorksheet, [translatedItem], { origin: -1 }); // 订单
                const childData = item.children || [];
                delete item.children;

                if (childData.length > 0) {
                    // 处理子记录字段
                    const translatedChildData = childData.map(childItem => {
                        const translatedChildItem = {};
                        for (const key in childItem) {
                            if (purchaseSummaryDetailFileds[key]) {
                                translatedChildItem[purchaseSummaryDetailFileds[key]] = childItem[key];
                            }
                        }
                        return translatedChildItem;
                    });
                    XLSXStyle.utils.sheet_add_json(jsonWorksheet, [{ "": TitleDetail }], { origin: -1 }); // 小标题
                    XLSXStyle.utils.sheet_add_json(jsonWorksheet, translatedChildData, { origin: -1 }); // 订单明细
                }
            });
            workBook.SheetNames.push(supplierName);
            workBook.Sheets[supplierName] = jsonWorksheet;
        }
    }
    return XLSXStyle.writeFile(workBook, workBookName);
};
export default ExportExcelUtity;

// // 使用示例
// const data = [
//     // ... 这里放入你提供的数据
// ];

// const fileName = "exported_data.xlsx";
// ExportExcelUtity(data, fileName, "NodeliveredPurchase");
