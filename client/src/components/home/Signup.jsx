import React, { useState } from "react";
import 'bootstrap/dist/css/bootstrap.min.css';
import { Button, Form, FormGroup, Input, Label } from "reactstrap";
import {Link} from "react-router-dom"
import "./Signup.css"
import axios from "axios"

const Signup = () => {
  const [firstname, setFirstName] = useState("")
  const [lastname, setLastName] = useState("")
  const [phoneNumber, setPhone] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [errorMessage, setMessage] = useState("")

  // axios config
  const axiosConfig = {
    headers: {
      "Content-type" : "application/json"
    }
  };

  // Register as a new user
  const Signup = async () => {
    try{
      if(!firstname || !lastname || !phoneNumber || !password || !confirmPassword) {
        return
      }
      if (password !== confirmPassword) {
        return
      }
      const {data} = await axios.post("/signup", 
      {firstName : firstname, lastName : lastname, phoneNumber : phoneNumber, password : password}, axiosConfig)

      setMessage("Registration Successful!")
    }catch(err) {
      setMessage("User with this phone number number already exists")
    }
  }



  return (
    <>
     <div className="register_form">
        <FormGroup>
          <Label for="firstname">First Name</Label>
          <Input id="firstname" name="firstname" type="text"  value={firstname} onChange={(e) => setFirstName(e.target.value)} />
          <Label for="firstname">Last Name</Label>
          <Input id="lastname" name="lastname" type="text"  value={lastname} onChange={(e) => setLastName(e.target.value)} />
        </FormGroup>
        <FormGroup>
          <Label for="phonenumber">Phone Number</Label>
          <Input id="phonenumber" name="phonenumber" type="number" value={phoneNumber} onChange={(e) => setPhone(e.target.value)} />
        </FormGroup>
        <FormGroup>
          <Label for="password">Password</Label>
          <Input id="password" name="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </FormGroup>
        <FormGroup>
          <Label for="confirm_password">Confirm Password</Label>
          <Input id="confirm_password" name="confirm_password" type="password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} />
        </FormGroup>
        <div className="error_div">
          <p id="error_message">{errorMessage}</p>
        </div>
        <Button color="info" onClick={Signup}>Register</Button>
      </div>
      <div className="info_text">
      <Link to={"/"} id="text">
          Already have an account? Login here.
        </Link>
      </div>
    </>
  );
};

export default Signup;
