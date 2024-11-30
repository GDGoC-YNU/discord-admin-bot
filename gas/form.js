function onFormSubmit(e) {
    const formResponses = e.response.getItemResponses();
    const targetQuestionTitle = "管理用ID"; // 特定の質問タイトル
    let userID = null;
    for (var i = 0; i < formResponses.length; i++) {
        var itemResponse = formResponses[i];
        if (itemResponse.getItem().getTitle().startsWith(targetQuestionTitle)) {
            userID = itemResponse.getResponse(); // 回答を取得
            break;
        }
    }
    if (!userID) {
        Logger.log("Could not identify the user ID.");
        return;
    }
    Logger.log("Submitted user ID: " + userID);

    try {
        let endpointUrl = "https://bot.gdgoc-ynu.shion.pro/api/initial/form/submit"
        var payload = {"user_id": userID}
        var options = {
            method: "post",
            contentType: "application/json",
            payload: JSON.stringify(payload)
        };
        var resp = UrlFetchApp.fetch(endpointUrl, options);
        if (resp.getResponseCode() === 200) {
            Logger.log("Successfully notified to the server.");
        }else {
            Logger.log("Failed to send the form data to the server, response code: " + resp.getResponseCode());
        }
    } catch (error) {
        Logger.log("Caught error: " + error.message);
    }
}
