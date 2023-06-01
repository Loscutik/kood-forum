function header() {
  if (document.body.scrollTop == 0 || document.documentElement.scrollTop == 0) {
    document.getElementById("header").style.marginTop = "0px";
  }
  if (document.body.scrollTop > 50 || document.documentElement.scrollTop > 50) {
    document.getElementById("header").style.marginTop = "-57px";
  }
}
function header2() {
  document.getElementById("header").style.marginTop = "0px";
}
function header3() {
  if (document.body.scrollTop > 50 || document.documentElement.scrollTop > 50) {
    document.getElementById("header").style.marginTop = "-57px";
  }
}

function signIn() {
  document.getElementById("darkness").style.display = "block";
  document.getElementById("signinform").style.display = "block";
  document.getElementById("signupform").style.display = "none";
}
function signUp() {
  document.getElementById("darkness").style.display = "block";
  document.getElementById("signupform").style.display = "block";
  document.getElementById("signinform").style.display = "none";

}

function checkFormSignup() {
  const submitButton = document.querySelector("#signup_submit");
  const email = document.querySelector("#email")
  const name = document.querySelector("#name-up")
  const password = document.getElementById("password-up");
  const confirmPassword = document.getElementById("confirm_password");
  const warning = document.querySelector("#warning-up");

  submitButton.addEventListener("click", (event) => {
    if (email.value == null || name.value == null || password.value == null || confirmPassword.value == null || email.value == "" || name.value == "" || password.value == "" || confirmPassword.value == "") {
      warning.innerHTML = "error: fill all fields";
      warning.style.display = "block";
    } else if (password.value != confirmPassword.value) {
      warning.innerHTML = "error: passwords do not match";
      warning.style.display = "block";
    } else {
      // send data by post
      // form data for post
      let data = {
        email: email.value,
        name: name.value,
        password: password.value,
      };
      sendPost(data, "/signup", warning, goIfSuccess);
    }
  });
}

function checkFormSignin() {
  const submitButton = document.querySelector("#signin_submit");
  const name = document.querySelector("#name-in")
  const password = document.getElementById("password-in");
  const warning = document.querySelector("#warning-in");

  submitButton.addEventListener("click", (event) => {
    if (name.value == null || password.value == null || name.value == "" || password.value == "") {
      warning.innerHTML = "error: fill all fields";
      warning.style.display = "block";
    } else {
      // send data by post
      // form data for post
      let data = {
        name: name.value,
        password: password.value,
      };

      sendPost(data, "/login", warning, goIfSuccess)
    }
  });
}

function changingForm(formID, warningID) {
  document.getElementById(formID).addEventListener("input", (event) => {
    document.getElementById(warningID).style.display = "none";
    document.getElementById(warningID).style.color = "rgb(239, 77, 93)";
  })
}

function changingFormSignup() {
  const form = document.querySelector("#signup_form");
  form.addEventListener("change", (event) => {
    document.getElementById("warning-up").style.display = "none";
  })
}

function changingFormSignin() {
  const form = document.querySelector("#signin_form");
  form.addEventListener("change", (event) => {
    document.getElementById("warning-in").style.display = "none";
  })
}

function handleLike(id) {
  // needed : "messageType"("posts_likes", "comments_likes") "messageID"(#)  "like"(bool) 
  const clickedElement = document.getElementById(id);
  let messageType = clickedElement.getAttribute("messageType");
  let messageID = clickedElement.getAttribute("messageID");
  const labelLike = document.getElementById(messageID + "-" + messageType + "-true-n");
  const labelDislike = document.getElementById(messageID + "-" + messageType + "-false-n");
  // create a request with JSON data
  let data = {
    messageType: messageType,
    messageID: messageID,
    like: clickedElement.getAttribute("like"),
  };
  const headers = new Headers();
  headers.append('Content-Type', 'application/json');

  fetch("/liking", {
    method: "POST",
    headers: headers,
    credentials: "same-origin",
    redirect: "follow",
    body: JSON.stringify(data)
  }).then(res => {
    if (!res.ok) {
      throw new Error(`HTTP error! Status: ${res.status}`);
    }
    return res.json();
  })
    .then(likes => {
      labelLike.innerHTML = likes["like"];
      labelDislike.innerHTML = likes["dislike"];
    });
}

function imageDownload(id) {
  // needed : "messageType"("p", "c") "messageID"(#) 
  const clickedElement = document.getElementById(id);

  const imageFile = clickedElement.files?.[0];
  if (!imageFile) { return; }
  // create a request

  const formData = new FormData();

  formData.append("messageType", clickedElement.getAttribute("messageType"));
  formData.messageID("accountnum", clickedElement.getAttribute("messageID"));

  // HTML file input, chosen by user
  formData.append("imagefile", imageFile);


  // send the POST request to the server
  fetch("/imagedownload", {
    method: "POST",
    // headers: headers,
    credentials: "same-origin",
    redirect: "error",
    body: formData
  }).then((res) => {
    if (!res.ok) {
      throw new Error(`HTTP error! Status: ${res.status}`);
    }
    return res.text();
  }).then((fileUrl) => {
    const imageDiv = document.getElementById("images");
    const img = document.createElement('img');
    img.setAttribute("src", fileUrl);
    imageDiv.appendChild(img);
  });
}

function darkness() {
  document.getElementById("darkness").style.display = "none";
  document.getElementById("signinform").style.display = "none";
  document.getElementById("signupform").style.display = "none";
}

function openSidepanel() {
  if (document.getElementById("usersidepanel").style.display == "none") {
    document.getElementById("usersidepanel").style.display = "block";
  } else {
    document.getElementById("usersidepanel").style.display = "none";
  }
}

function openFilters() {
  if (document.getElementById("filterform").style.display == "none") {
    document.getElementById("filterform").style.display = "block";
    document.getElementById("filterslogo").style.display = "none";
    document.getElementById("closefilterslogo").style.display = "block";
  } else {
    document.getElementById("filterform").style.display = "none";
    document.getElementById("filterslogo").style.display = "block";
    document.getElementById("closefilterslogo").style.display = "none";
  }
}

function ShowPassword() {
  var x = document.getElementsByClassName("password");
  for (let i = 0; i < x.length; i++) {
    if (x[i].type === "password") {
      x[i].type = "text";
    } else {
      x[i].type = "password";
    }
  }
}

function checkFormSettings() {
  const email = document.getElementById("email")
  const submitEmail = document.getElementById("submit_email");
  const warningEmail = document.getElementById("warning_email");
  const password = document.getElementById("password");
  const confirmPassword = document.getElementById("confirm_password");
  const submitPassword = document.getElementById("submit_password");
  const warningPassword = document.getElementById("warning_password");

  submitEmail.addEventListener("click", event => {
    if (email.value == null || email.value == "") {
      warningEmail.innerHTML = "error: fill all fields";
      warningEmail.style.display = "block";
    } else {
      // send data by post
      // form data for post
      let data = {
        email: email.value,
      };

      sendPost(data, "/settings", warningEmail, (res => { }));
    }
  });

  submitPassword.addEventListener("click", event => {
    if (password.value == null || confirmPassword.value == null || password.value == "" || confirmPassword.value == "") {
      warningPassword.innerHTML = "error: fill all fields";
      warningPassword.style.display = "block";
    } else if (password.value != confirmPassword.value) {
      warningPassword.innerHTML = "error: passwords do not match";
      warningPassword.style.display = "block";
    } else {
      // send data by post
      // form data for post
      let data = {
        password: password.value,
      };

      sendPost(data, "/settings", warningPassword, (res => { }));
    }
  });
}

async function goIfSuccess(res) {
  if (res.status == 204) {
    window.location.href = res.headers.get("Location");
  }
}

const sendPost = async (data, url, warningElm, checkSpecialCase) => {
  // create a request with form-data
  const urlEncodedDataPairs = [];
  for (const [name, value] of Object.entries(data)) {
    urlEncodedDataPairs.push(`${encodeURIComponent(name)}=${encodeURIComponent(value)}`);
  }

  // Combine the pairs into a single string and replace all %-encoded spaces to
  // the '+' character; matches the behavior of browser form submissions.
  const urlEncodedData = urlEncodedDataPairs.join('&').replace(/%20/g, '+');
  const headers = new Headers();
  headers.append('Content-Type', 'application/x-www-form-urlencoded');

  // send the POST request to the server
  const res = await fetch(url, {
    method: "POST",
    headers: headers,
    credentials: "same-origin",
    redirect: "error",
    body: urlEncodedData
  })

  if (!res.ok) {
    const html = await res.text();
    document.querySelector("html").innerHTML = html;
    return;
  } else {
    checkSpecialCase(res);
    const text = await res.text();
    if (text.length != 0) {
      if (!text.startsWith("error:")) {
        warningElm.style.color = "rgb(57, 202, 62)"
      }
      warningElm.innerHTML = text;
      warningElm.style.display = "block";
    }
  }
}

function validatePost() {
  var x = document.forms["pform"]["theme"].value;
  var y = document.forms["pform"]["content"].value;
  if (x.trim() == "") {
    document.getElementById("PostTopic").style.border = "solid 2px";
    document.getElementById("PostTopic").style.borderColor = "rgb(232, 0, 0)";
    document.getElementById("PostTopic").style.borderRadius = "3px";
    document.getElementById("PostTopic").placeholder = "Please enter the topic!";
  }
  if (y.trim() == "") {
    document.getElementById("textarea_newpost").style.border = "solid 2px";
    document.getElementById("textarea_newpost").style.borderColor = "rgb(232, 0, 0)";
    document.getElementById("textarea_newpost").style.borderRadius = "3px";
    document.getElementById("textarea_newpost").placeholder = "Please enter the text!";
  }

  if (x.trim() == "" || y.trim() == "") {
    return false
  }
  else {
    return true
  }
}

function CheckValidatePost() {
  if (validatePost() == false) {
    document.getElementById("PostTopic").style.border = "solid 1px";
    document.getElementById("PostTopic").style.borderColor = "rgb(0, 0, 0)";
    document.getElementById("PostTopic").style.borderRadius = "3px";
    document.getElementById("PostTopic").placeholder = "Header";
    document.getElementById("textarea_newpost").style.border = "solid 1px";
    document.getElementById("textarea_newpost").style.borderColor = "rgb(0, 0, 0)";
    document.getElementById("textarea_newpost").style.borderRadius = "3px";
    document.getElementById("textarea_newpost").placeholder = "Enter your text here...";
  }
  if (CheckCheckBox() == false) {
    document.querySelectorAll('.categorylabel').forEach(el => {
      el.style.border = "solid rgb(255, 193, 47) 2px";
    })
    document.getElementById("choosetags").style.color = "rgb(0, 0, 0)";
  }
}

function Up() {
  if (validatePost() == false || CheckCheckBox() == false) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
}

document.addEventListener('DOMContentLoaded', function () {

  document.querySelectorAll('.categorylabel').forEach(el => {
    el.addEventListener("click", function (ev) {
      ev.target.classList.toggle("selected")
    })
  });

  let pform = document.getElementById("pform")
  if (pform) {
    document.getElementById("pform").addEventListener("submit", function (ev) {
      console.log(ev.target)
      ev.preventDefault();
      if (validatePost() && CheckCheckBox()) {
        ev.target.submit();
      }
    })
  }


}, false);

function CheckCheckBox() {
  var t = false;
  document.querySelectorAll('.categories').forEach(el => {
    if (el.checked == true) {
      t = true
    }
  })
  if (t == false) {
    document.querySelectorAll('.categorylabel').forEach(el => {
      el.style.border = "solid rgb(232, 0, 0) 2px";
    })
    document.getElementById("choosetags").style.color = "rgb(255, 0, 0)";
    return false
  }
  return true
}


function validateComment() {
  var z = document.forms["writecomment"]["content"].value;

  const inputImages = document.getElementById('image_uploads');

  const curFiles = inputImages.files;
  let isImagesValide = true;

  for (const file of curFiles) {
    if (!validFileType(file) || file.size > 2 * 1024 * 1024) {
      isImagesValide = false;
      break;
    }
  }

  if ((z.trim() == "" && curFiles.length === 0) || !isImagesValide) {
    document.getElementById("newcomment").style.border = "solid 2px";
    document.getElementById("newcomment").style.borderColor = "rgb(232, 0, 0)";
    document.getElementById("newcomment").style.borderRadius = "3px";
    document.getElementById("newcomment").placeholder = "Can't submit an empty comment";
    return false
  } else {
    return true
  }
}

function checkComment() {
  if (document.getElementById("newcomment").style.borderColor = "rgb(232, 0, 0)") {
    document.getElementById("newcomment").style.border = "solid 1px";
    document.getElementById("newcomment").style.borderColor = "rgb(0, 0, 0)";
    document.getElementById("newcomment").style.borderRadius = "3px";
    document.getElementById("newcomment").placeholder = "Write your comment...";
  }
}

function choseImage() {

  const inputs = document.querySelectorAll('.image_uploads');

  inputs.forEach((input) => { input.addEventListener('change', updateImageDisplay) });
}

function updateImageDisplay(event) {
  const targetElmID = event.target.id;
  const input = document.getElementById(targetElmID);
  const preview = document.getElementById(targetElmID.replace('image_uploads', 'preview'));
  while (preview.firstChild) {
    preview.removeChild(preview.firstChild);
  }

  const curFiles = input.files;
  if (curFiles.length === 0) {
    const para = document.createElement('p');
    para.textContent = 'No files currently selected for upload';
    preview.appendChild(para);
  } else {
    const list = document.createElement('ol');
    preview.appendChild(list);

    for (const file of curFiles) {
      const listItem = document.createElement('li');
      const para = document.createElement('p');
      if (validFileType(file)) {
        if (file.size > 2 * 1024 * 1024) {
          para.textContent = 'file is too big';
        } else {
          para.textContent = `File name ${file.name}, file size ${returnFileSize(file.size)}.`;
          const image = document.createElement('img');
          image.src = URL.createObjectURL(file);
          listItem.appendChild(image);
        }
        listItem.appendChild(para);
      } else {
        para.textContent = `File name ${file.name}: Not a valid file type. Update your selection.`;
        listItem.appendChild(para);
      }

      list.appendChild(listItem);
    }
  }
}

const fileTypes = [
  "image/bmp",
  "image/gif",
  "image/jpeg",
  "image/jpg",
  "image/png",
  "image/svg+xml",
];

function validFileType(file) {
  return fileTypes.includes(file.type);
}

function returnFileSize(number) {
  if (number < 1024) {
    return `${number} bytes`;
  } else if (number >= 1024 && number < 1048576) {
    return `${(number / 1024).toFixed(1)} KB`;
  } else if (number >= 1048576) {
    return `${(number / 1048576).toFixed(1)} MB`;
  }
}

