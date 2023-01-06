import React from "react";
import { Button, Form, FormGroup, Input, Label } from "reactstrap";
import {Link, useNavigate} from "react-router-dom"
import "./Login.css";
import { useState, useEffect } from "react";
import axios from "axios"

const Login = () => {
  let navigate = useNavigate()
  const [phone, setPhone] = useState("")
  const [password, setPassword] = useState("")

    // axios config
    const axiosConfig = {
      headers: {
        "Content-type" : "application/json"
      }
    };

   // login User
   const Login = async () => {
    try{
      if(!phone || !password) {
        return
      }

      const {data} = await axios.post("/login", 
      {phoneNumber : phone, password: password}, axiosConfig)

      localStorage.setItem("token", data.token);
      localStorage.setItem("firstname", data.firstname);
      localStorage.setItem("lastname", data.lastname);
      localStorage.setItem("phone", data.phone);
      navigate("/dashboard")
    }catch(err) {
      console.log(err)
    }
  }

  useEffect(() => {
    // If user is alraedy loggedin, then redirect to dashboard
    if (localStorage.getItem("token") != null) {
      navigate("/dashboard")
    }
  }, [])
  
  return (
    <div>
      {" "}
      <div className="login_form">
        <FormGroup>
          <Label for="phonenumber">Phone Number</Label>
          <Input id="phonenumber" name="phonenumber" type="number" value={phone}  onChange={(e) => setPhone(e.target.value)}/>
        </FormGroup>
        <FormGroup>
          <Label for="password">Password</Label>
          <Input id="password" name="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)}/>
        </FormGroup>
        <Button color="info" onClick={Login}>Login</Button>
      </div>
      <div className="info_text">
        <Link to={"/register"} id="text">
          Don't have an account? Register now.
        </Link>
      </div>
    </div>
  );
};

export default Login;
