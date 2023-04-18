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
      });
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
      });
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