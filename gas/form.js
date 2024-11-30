function onFormSubmit(e) {
  const formResponses = e.response.getItemResponses();
  const targetQuestionTitle = "管理用ID";
  let userID = null;
  for (let i = 0; i < formResponses.length; i++) {
    const itemResponse = formResponses[i];
    if (itemResponse.getItem().getTitle().startsWith(targetQuestionTitle)) {
      userID = itemResponse.getResponse();
      break;
    }
  }
  if (!userID) {
    throw new Error("Could not identify the user ID.");
  }
  Logger.log("Submitted user ID: " + userID);

  try {
    let endpointUrl = "https://bot.gdgoc-ynu.shion.pro/api/initial/form/submit"
    const payload = {"user_id": userID};
    const options = {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify(payload)
    };
    const resp = UrlFetchApp.fetch(endpointUrl, options);
    if (resp.getResponseCode() === 200) {
      Logger.log("Successfully notified to the server.");
    } else {
      throw new Error("Failed to send the form data to the server, response code: " + resp.getResponseCode());
    }
  } catch (error) {
    throw new Error("Caught error: " + error.message);
  }
}
