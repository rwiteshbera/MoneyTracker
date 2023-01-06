import React from "react";
import { useState, useEffect } from "react";
import "./Dashboard.css";
import { Button, Container, Input } from "reactstrap";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const Dashboard = () => {
  let navigate = useNavigate();

  const [firstname, setFirstName] = useState("");
  const [lastname, setLastName] = useState("");
  const [amount, setAmount] = useState("");

  let authToken = localStorage.getItem("token");
  let phone = localStorage.getItem("phone")

  // axios config
  const axiosConfig = {
    headers: {
      "Content-type": "application/json",
      Authorization: authToken,
    },
  };

  // create new transaction
  const createTransaction = async () => {
    try {
        if(!amount || !phone) {
            return
        }
      const { data } = await axios.post(
        "/create",
        { createdBy: phone, amount: amount },
        axiosConfig
      );
      console.log(data);
    } catch (e) {
      console.log(e);
    }
  };

  // Logout the account
  const Logout = () => {
    localStorage.clear();
    navigate("/"); // redirecting to login page after logout
  };

  useEffect(() => {
    setFirstName(localStorage.getItem("firstname"));
    setLastName(localStorage.getItem("lastname"));
  }, []);

  return (
    <>
      <Container>
        <h4 className="username_heading">
          Hi, {firstname} {lastname}
        </h4>
        <div className="parent_box">
          <div className="amount_box">
            <Input
              type="number"
              value={amount}
              id="input_amount"
              onChange={(e) => setAmount(e.target.value)}
            />
            <Button color="success" onClick={createTransaction}>
              Create Transacion
            </Button>
            <Button color="warning" onClick={Logout} id="logout_btn">
              Logout
            </Button>
          </div>
        </div>
      </Container>
    </>
  );
};

export default Dashboard;
