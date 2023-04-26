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

function checkFormSignup(){ 
  const submitButton=document.querySelector("#signup_submit");
  const email=document.querySelector("#email") 
  const name= document.querySelector("#name-up")
  const password=document.getElementById("password-up");
  const confirmPassword=document.getElementById("confirm_password");
  const warning=document.querySelector("#warning-up");

  submitButton.addEventListener("click", (event) => {
    if (email.value==null ||name.value==null ||password.value==null ||confirmPassword.value==null|| email.value=="" || name.value=="" ||password.value=="" ||confirmPassword.value==""){
      warning.innerHTML="fill all fields";
      warning.style.display="block";
    }else if (password.value!=confirmPassword.value){
      warning.innerHTML="passwords do not match";
      warning.style.display="block";
    }else{
      // send data by post
      // form data for post
      let data = { 
        email: email.value,
        name: name.value,
        password: password.value,
      };
      
      sendPost(data, "/signup", warning, goIfSuccess);
      /*
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

      fetch("/signup", {
      method: "POST",
      headers: headers,
      credentials: "same-origin",
      redirect: "follow", 
      body: urlEncodedData
      }).then((res) => {
        if (res.status==204){
          console.log ("red=",res.headers.get("Location"));
          window.location.href =res.headers.get("Location");
          return ""
        }
        if (!res.ok) {
          throw new Error(`HTTP error! Status: ${res.status}`);
        }
        
        return res.text(); 
      })
      .then((text) =>{
        if (text.length!=0){
          warning.innerHTML=text;
          warning.style.display="block";
        }
      });*/
    }
  });
}

function checkFormSignin(){
  const submitButton=document.querySelector("#signin_submit");
  const name= document.querySelector("#name-in")
  const password=document.getElementById("password-in");
  const warning=document.querySelector("#warning-in");

  submitButton.addEventListener("click", (event) => {
    if (name.value==null ||password.value==null || name.value=="" ||password.value=="" ){
      warning.innerHTML="fill all fields";
      warning.style.display="block";
    }else{
      // send data by post
      // form data for post
      let data = { 
        name: name.value,
        password: password.value,
      };

      sendPost(data, "/login", warning, goIfSuccess)
      
      /*
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

      fetch("/login", {
      method: "POST",
      headers: headers,
      credentials: "same-origin",
      redirect: "follow", 
      body: urlEncodedData
      }).then((res) => {
        if (res.status==204){
          console.log ("red=",res.headers.get("Location"));
          window.location.href =res.headers.get("Location");
          return ""
        }
        if (!res.ok) {
          throw new Error(`HTTP error! Status: ${res.status}`);
        }
        
        return res.text(); 
      })
      .then((text) =>{
        if (text.length!=0){
          warning.innerHTML=text;
          warning.style.display="block";
        }
      });*/
    }
  });
}

function changingFormSignup() {
  const form=document.querySelector("#signup_form");
  form.addEventListener("change", (event) =>{
    document.getElementById("warning-up").style.display="none";
  })
}

function changingFormSignin() {
  const form=document.querySelector("#signin_form");
  form.addEventListener("change", (event) =>{
    document.getElementById("warning-in").style.display="none";
  })
}

function handleLike(id){
  // needed : "messageType"("posts_likes", "comments_likes") "messageID"(#)  "like"(bool) 
  const clickedElement = document.getElementById(id);
  var messageType = clickedElement.getAttribute("messageType");
  var messageID = clickedElement.getAttribute("messageID");
  const labelLike = document.getElementById(messageID+"-"+messageType+"-true-n");
  const labelDislike = document.getElementById(messageID+"-"+messageType+"-false-n");
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
  .then(likes =>{
    labelLike.innerHTML=likes["like"];
    labelDislike.innerHTML=likes["dislike"];
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

function checkFormSettings(){ 
  const email=document.getElementById("email") 
  const submitEmail=document.getElementById("submit_email");
  const warningEmail=document.getElementById("warning_email");
  const password=document.getElementById("password");
  const confirmPassword=document.getElementById("confirm_password");
  const submitPassword=document.getElementById("submit_password");
  const warningPassword=document.getElementById("warning_password");
  
  submitEmail.addEventListener("click", event => {
    if (email.value==null || email.value==""){
      warningEmail.innerHTML="fill all fields";
      warningEmail.style.display="block";
    }else{
      // send data by post
      // form data for post
      let data = { 
        email: email.value,
      };
      
      sendPost(data, "/settings", warningEmail, (res=>{})); 
    }
  });

  submitPassword.addEventListener("click", event => {
    if (password.value==null ||confirmPassword.value==null || password.value=="" || confirmPassword.value==""){
      warningPassword.innerHTML="fill all fields";
      warningPassword.style.display="block";
    }else if (password.value!=confirmPassword.value){
      warningPassword.innerHTML="passwords do not match";
      warningPassword.style.display="block";
    }else{
      // send data by post
      // form data for post
      let data = { 
        password: password.value,
      };
      
      sendPost(data, "/settings", warningPassword, (res=>{})); 
    }
  });
}

async function goIfSuccess(res){
  if (res.status==204){
    window.location.href =res.headers.get("Location");
  }
}

const sendPost = async  (data, url, warningElm, checkSpecialCase)=>{
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
   const res= await fetch(url, {
     method: "POST",
     headers: headers,
     credentials: "same-origin",
     redirect: "error", 
     body: urlEncodedData
   })

    if (!res.ok){
      const html=  await res.text();
      document.querySelector("html").innerHTML=html;
      return; 
    }else{
      checkSpecialCase(res);
      const text = await res.text();
      if (text.length!=0){
        warningElm.innerHTML=text;
        warningElm.style.display="block";
      }
    }
}

function validatePost() {
  var x = document.forms["pform"]["theme"].value;
  var y = document.forms["pform"]["content"].value;
  if (x == "") {
    document.getElementById("PostTopic").style.border = "solid 2px";
    document.getElementById("PostTopic").style.borderColor = "rgb(232, 0, 0)";
    document.getElementById("PostTopic").style.borderRadius = "3px";
    document.getElementById("PostTopic").placeholder = "Please enter the topic!";
  }
  if (y == "") {
    document.getElementById("textarea_newpost").style.border = "solid 2px";
    document.getElementById("textarea_newpost").style.borderColor = "rgb(232, 0, 0)";
    document.getElementById("textarea_newpost").style.borderRadius = "3px";
    document.getElementById("textarea_newpost").placeholder = "Please enter the text!";
  }
  if (x == "" || y == "") {
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
}

function Up() {
  if (validatePost() == false) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
}

document.addEventListener('DOMContentLoaded', function () {
  document.querySelectorAll('.categorylabel').forEach(el => {
    el.addEventListener("click", function(ev){
      ev.target.classList.toggle("selected")
    })
  });
}, false);