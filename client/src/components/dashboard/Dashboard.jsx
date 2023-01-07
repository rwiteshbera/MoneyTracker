import React from "react";
import { useState, useEffect } from "react";
import "./Dashboard.css";
import {
  Button,
  Container,
  Input,
  Accordion,
  AccordionBody,
  AccordionHeader,
  AccordionItem,
  Modal,
  ModalHeader,
  ModalBody,
  ModalFooter,
} from "reactstrap";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const Dashboard = () => {
  let navigate = useNavigate();

  const [firstname, setFirstName] = useState("");
  const [lastname, setLastName] = useState("");
  const [amount, setAmount] = useState("");
  const [transactionName, setTransactionName] = useState("");
  const [transactionList, setTransactionList] = useState([]);
  const [searchPhoneInput, setSeachPhoneInput] = useState("");
  const [searchMemberResult, setSearchMemberResult] = useState("");
  const [addButtonTextInModal, setAddButtonTextInModal] = useState("Add")
  const [allMembers, setAllMembers] = useState([])

  const [open, setOpen] = useState("1");
  const toggle = (id) => {
    if (open === id) {
      setOpen();
    } else {
      setOpen(id);
    }
  };

  const [modal, setModal] = useState(false);
  const [currentModalId, setCurrentModalId] = useState("")
  const toggleModal = () => {
    setModal(!modal)
    setSeachPhoneInput("")
    setSearchMemberResult("")
    setAddButtonTextInModal("Add")
  };

  let authToken = localStorage.getItem("token");
  let phone = localStorage.getItem("phone");

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
      if (!amount || !phone) {
        return;
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

  // Get a list of all transactions created by
  const getTransactionsList = () => {
    try {
      axios.get("/get", axiosConfig).then((res) => {
        setTransactionList(res.data?.message.reverse());
      });
    } catch (error) {
      console.log(error);
    }
  };

  // Search memeber by phone number
  const searchUser = async () => {
    try {
      const { data } = await axios.post(
        "/search",
        { phoneNumber: searchPhoneInput },
        axiosConfig
      );
      setSearchMemberResult(data.user);
      console.log(data.user);
    } catch (error) {
      console.log(error);
    }
  };

  // Add memeber by phone number
  const addMemeber = async (currentModalId, phone) => {
    try {
      const { data } = await axios.post(
        "/add_member",
        { phone_number: phone, transaction_id: currentModalId },
        axiosConfig
      );
      setAddButtonTextInModal("Added")
      console.log(data);
    } catch (error) {
      console.log(error);
    }
  };


  // Show members
  const showMemebers = async (currentModalId, phone) => {
    try {
      const { data } = await axios.get(
        "/show",
        axiosConfig
      );
      setAllMembers(data.message)
    } catch (error) {
      console.log(error);
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
    getTransactionsList();
    showMemebers();
  }, []);

  useEffect(() => {
    getTransactionsList();
  }, [transactionList]);

  return (
    <>
      <Container>
        <h4 className="username_heading">
          Hi, {firstname} {lastname}
        </h4>
        <div className="parent_box">
          <div className="transaction_box">
            <Input
              type="text"
              value={transactionName}
              id="transaction_name_input"
              placeholder="Enter transaction name"
              onChange={(e) => setTransactionName(e.target.value)}
            />
            <Input
              type="number"
              value={amount}
              id="input_amount"
              placeholder="Transaction amount"
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
        <div className="transaction_list">
          <Accordion open={open} toggle={toggle}>
            {transactionList.map(({ id, amount }, key) => {
              return (
                <>
                  <div key={key}>
                    <AccordionHeader targetId="1">
                      Transaction Id: {id} | Amount: {amount}
                    </AccordionHeader>
                    <AccordionBody
                      accordionId="1"
                      className="tranasaction_accordion_body"
                    >
                      <strong>{amount}</strong>
                      <Button className="addmember_btn" onClick={() => {toggleModal();setCurrentModalId(id)}}>
                        Add Member
                      </Button>

                      {
                        allMembers.map(({phone_number, amount_to_be_paid}, key) => {
                          return (
                              <>
                                  <div>{phone_number}</div>
                                  <div>{amount_to_be_paid}</div>
                              </>
                          );
                        })
            }
                    </AccordionBody>
                  </div>
                </>
              );
            })}
          </Accordion>
        </div>
      </Container>

      <Modal
        isOpen={modal}
        fade={false}
        toggle={toggleModal}
        className="search_modal"
      >
        <ModalHeader toggle={toggleModal}>Add Member</ModalHeader>
        <ModalBody className="search_modal_body">
          <div className="search_modal_cointainer">
            <Input
              type="text"
              placeholder="Search phone"
              value={searchPhoneInput}
              onChange={(e) => setSeachPhoneInput(e.target.value)}
            />
            {searchPhoneInput.length === 10 && (
              <Button onClick={searchUser}>Search</Button>
            )}
          </div>
          {searchMemberResult.firstName !== undefined && (
            <div className="search_member_result">
              {searchMemberResult.firstName + " " + searchMemberResult.lastName}{" "}
              <Button
                className="modal_add_btn"
                color="success"
                onClick={() => addMemeber(currentModalId, searchPhoneInput)}
              >
                {addButtonTextInModal}
              </Button>
            </div>
          )}
        </ModalBody>
        <ModalFooter>
          <Button color="secondary" onClick={toggleModal}>
            Cancel
          </Button>
        </ModalFooter>
      </Modal>
    </>
  );
};

export default Dashboard;
