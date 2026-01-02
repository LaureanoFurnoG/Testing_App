import "./style.css";
import React, { useEffect, useState } from "react";
import axiosInstance from "../../axios";
import { Divider, Table, Input } from "antd";
import type { TableColumnsType } from "antd";
import { DeleteOutlined } from "@ant-design/icons";
import AddUserDrawer from "../AddUser-Drawer/AddUser-Drawer";
import { useParams } from "react-router-dom";

const { Search } = Input;

interface DataType {
  key: React.Key;
  User: React.ReactNode;
}

export default function TableUsers_Group() {
  const { groupId } = useParams();
  const [data, setData] = useState<DataType[]>([]);

  const deleteUser = (key: React.Key) => {
    console.log("Delete:", key);
  };

  const getMembers = async () => {
    try {
      const response = await axiosInstance.get(
        `/api/group/showAllUsersGroup/${groupId}`
      );

      const formattedData = response.data.groupMembers?.map(
        (user: any) => ({
          key: user.id,
          User: (
            <div className="profileData">
              <img src={user.avatar || ""} className="ProfileImage" alt="" />
              <p>{user.name}</p>
            </div>
          ),
        })
      );
      setData(formattedData);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    getMembers();
  }, [groupId]);

  const columns: TableColumnsType<DataType> = [
    {
      title: "User",
      dataIndex: "User",
    },
    {
      title: "Delete",
      render: (_, record) => (
        <div className="deleteCell">
          <DeleteOutlined
            className="deleteUser-icon"
            onClick={() => deleteUser(record.key)}
          />
        </div>
      ),
    },
  ];

  return (
    <div>
      <div className="headerUsers">
        <span className="headerTitle">User</span>

        <div className="headerActions">
          <AddUserDrawer />
          <Search
            placeholder="Search user by Name"
            className="searchUser"
          />
        </div>
      </div>

      <Divider className="divider" />

      <Table
        columns={columns}
        dataSource={data}
        showHeader={false}
      />
    </div>
  );
}
