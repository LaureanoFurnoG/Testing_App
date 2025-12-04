import "./style.css";
import React from "react";
import axiosInstance from "../../axios";
import { Divider, Table, type GetProps } from "antd";
import { Input, TableColumnsType } from "antd";
import { DeleteOutlined } from "@ant-design/icons";

const { Search } = Input;

type SearchProps = GetProps<typeof Input.Search>;

interface DataType {
  key: React.Key;
  User: any;
  delete: any;
}

interface TableUsers_GroupProps {
  URL: string;
}

interface TableUsers_GroupState {
  data: DataType[];
}

export default class TableUsers_Group extends React.Component<
  TableUsers_GroupProps,
  TableUsers_GroupState
> {
  constructor(props: TableUsers_GroupProps) {
    super(props);

    this.state = {
      data: [
        {
          key: "1",
          User: (
            <div className="profileData">
              <img src="" className="ProfileImage" alt="" />
              <p>Name</p>
            </div>
          ),
          delete: true,
        },
      ],
    };
  }

  onSearch: SearchProps["onSearch"] = (value, _e, info) => {
    console.log(info?.source, value);
  };

  deleteUser = (key: React.Key): void => {
    console.log("Delete:", key);
  };

  render() {
    const columns: TableColumnsType<DataType> = [
      {
        title: "User",
        dataIndex: "User",
      },
      {
        title: "Delete",
        dataIndex: "delete",
        render: (_, record) => (
          <div className="deleteCell">
            <DeleteOutlined
              className="deleteUser-icon"
              onClick={() => this.deleteUser(record.key)}
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
            <button className="addMemberBtn">Add Member</button>

            <Search
              placeholder="Search user"
              onSearch={this.onSearch}
              className="searchUser"
            />
          </div>
        </div>

        <Divider className="divider" />

        <Table<DataType>
          columns={columns}
          dataSource={this.state.data}
          showHeader={false}
        />
      </div>
    );
  }
}
