import * as XLSXStyle from "xlsx-js-style";
import _, { pull } from "lodash";
import { PurchaseSummaryDetailFileds, PurchaseSummaryFileds } from "../model/purchaseSummary";

/**
 * 2023-09-21 应廖要求由原一个公司一个工作表的形式变为订单一个工作表和明细一个工作表
 */
const Title = {
    v: "采购订单",
    t: "s",
    s: { font: { bold: true, sz: 16, name: "宋体" }, alignment: { vertical: "center", horizontal: "center" }, fill: { fgColor: { rgb: "87CEEB" } } }
};

const TitleDetail = {
    v: "订单明细",
    t: "s",
    s: { font: { bold: true, sz: 16, name: "宋体" }, alignment: { vertical: "center", horizontal: "center" }, fill: { fgColor: { rgb: "ffff00" } } }
};

const ExportExcelUtity = (res, workBookName, type) => {
    let workBook = XLSXStyle.utils.book_new();

    if (type == "NodeliveredPurchase") {
        let purchaseSummaryLength = 0; // 订单单元格宽度
        let purchaseSUmmaryDeatilLength = 0; // 订单明细单元格宽度
        let SupplierNameGroups = _.groupBy(res, "SupplierName");
        let jsonWorksheet = XLSXStyle.utils.book_new(); // 订单表
        let jsonWorksheetDetail = XLSXStyle.utils.book_new(); // 订单明细表
        // 标题
        XLSXStyle.utils.sheet_add_aoa(jsonWorksheet, [[Title]]);
        XLSXStyle.utils.sheet_add_aoa(jsonWorksheetDetail, [[TitleDetail]]);
        //表头
        XLSXStyle.utils.sheet_add_json(jsonWorksheet, [PurchaseSummaryFileds], { skipHeader: true, origin: -1 });
        XLSXStyle.utils.sheet_add_json(jsonWorksheetDetail, [PurchaseSummaryDetailFileds], { skipHeader: true, origin: -1 });
        for (let supplierName in SupplierNameGroups) {
            let supplierData = SupplierNameGroups[supplierName];
            supplierData.forEach(item => {
                purchaseSummaryLength = Object.keys(item).length;
                // 订单明细;
                const childData = item.children || [];
                purchaseSUmmaryDeatilLength = Object.keys(childData[0]).length;
                delete item.children;
                if (childData.length > 0) {
                    // 写入明细
                    XLSXStyle.utils.sheet_add_json(jsonWorksheetDetail, childData, { skipHeader: true, origin: -1 }); // 订单明细
                }
            });
            // 写入订单
            XLSXStyle.utils.sheet_add_json(jsonWorksheet, supplierData, { skipHeader: true, origin: -1 }); // 订单
        }
        // 表头样式
        for (const key in jsonWorksheet) {
            if (Number(key.slice(1)) == 2) {
                jsonWorksheet[key].s = {
                    font: { bold: true, sz: 13, name: "宋体" }
                };
            }
        }
        for (const key in jsonWorksheetDetail) {
            if (Number(key.slice(1) == 2)) {
                jsonWorksheetDetail[key].s = {
                    font: { bold: true, sz: 13, name: "宋体" }
                };
            }
        }
        // 单元格宽度
        let cols = [];
        let col = {};
        for (let i = 1; i <= purchaseSUmmaryDeatilLength; i++) {
            if (i == purchaseSummaryLength) {
                jsonWorksheet["!cols"] = cols;
            }
            col.wpx = 100;
            cols.push(col);
        }
        jsonWorksheetDetail["!cols"] = cols;
        // console.log(cols); // len 12
        // console.log(purchaseSummaryLength); //11
        // console.log(purchaseSUmmaryDeatilLength); //12
        // 合并单元格
        const merge1 = [{ s: { r: 0, c: 0 }, e: { r: 0, c: purchaseSummaryLength - 2 } }];
        const merge2 = [{ s: { r: 0, c: 0 }, e: { r: 0, c: purchaseSUmmaryDeatilLength - 1 } }];
        jsonWorksheet["!merges"] = merge1;
        jsonWorksheetDetail["!merges"] = merge2;
        // 生成Excel表格
        workBook.SheetNames.push(Title.v);
        workBook.SheetNames.push(TitleDetail.v);
        workBook.Sheets[Title.v] = jsonWorksheet;
        workBook.Sheets[TitleDetail.v] = jsonWorksheetDetail;

        // for (let supplierName in SupplierNameGroups) {
        //     let supplierData = SupplierNameGroups[supplierName];
        //     let jsonWorksheet = XLSXStyle.utils.book_new();
        //     // 添加表头
        //     const header = {
        //         v: supplierName,
        //         t: "s",
        //         s: { font: { bold: true, sz: 18, name: "宋体" }, alignment: { vertical: "center", horizontal: "center" } }
        //     };

        //     XLSXStyle.utils.sheet_add_aoa(jsonWorksheet, [[header]]);
        //     const merge = [{ s: { r: 0, c: 0 }, e: { r: 0, c: 11 } }];
        //     jsonWorksheet["!merges"] = merge;
        //     // 遍历主记录和子记录
        //     supplierData.forEach(item => {
        //         // 主记录字段中文
        //         const translatedItem = {};
        //         // 处理主记录字段
        //         for (const key in item) {
        //             if (purchaseSummaryFileds[key]) {
        //                 translatedItem[purchaseSummaryFileds[key]] = item[key];
        //             }
        //         }
        //         XLSXStyle.utils.sheet_add_json(jsonWorksheet, [{ "": Title }], { origin: -1 }); // 小标题
        //         XLSXStyle.utils.sheet_add_json(jsonWorksheet, [translatedItem], { origin: -1 }); // 订单
        //         const childData = item.children || [];
        //         delete item.children;

        //         if (childData.length > 0) {
        //             // 处理子记录字段
        //             const translatedChildData = childData.map(childItem => {
        //                 const translatedChildItem = {};
        //                 for (const key in childItem) {
        //                     if (purchaseSummaryDetailFileds[key]) {
        //                         translatedChildItem[purchaseSummaryDetailFileds[key]] = childItem[key];
        //                     }
        //                 }
        //                 return translatedChildItem;
        //             });
        //             XLSXStyle.utils.sheet_add_json(jsonWorksheet, [{ "": TitleDetail }], { origin: -1 }); // 小标题
        //             XLSXStyle.utils.sheet_add_json(jsonWorksheet, translatedChildData, { origin: -1 }); // 订单明细
        //         }
        //     });
        //     workBook.SheetNames.push(supplierName);
        //     workBook.Sheets[supplierName] = jsonWorksheet;
        // }
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
