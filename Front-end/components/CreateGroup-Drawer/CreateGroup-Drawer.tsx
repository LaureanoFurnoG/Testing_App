import React, { useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Col, Drawer, Form, Input, Row, Space } from 'antd';
import axiosInstance from '../../axios';
import { useGroups } from '../../context/GroupsContext';

const CreateGroupDrawer: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();
  const { refreshGroups } = useGroups();

  const showDrawer = () => {
    setOpen(true);
  };

  const onClose = () => {
    setOpen(false);
  };

  const createGroup = async () => {
    try {
      const values = await form.validateFields();
      await axiosInstance.post("/api/group/createGroup", values);
      refreshGroups();
      onClose();
      form.resetFields();
    } catch (error) {
      console.log("Validation error:", error);
    }
  };


  return (
    <>
      <Button type="primary" className='addMemberBtn' onClick={showDrawer} icon={<PlusOutlined />}>
        Create Group
      </Button>

      <Drawer
        title="Create Group"
        onClose={onClose}
        open={open}
        styles={{
          body: {
            paddingBottom: 80,
          },
        }}
        extra={
          <Space>
            <Button className='btn-border' onClick={onClose}>Cancel</Button>
            <Button className='addMemberBtn' onClick={createGroup} type="primary">
              Create Group
            </Button>
          </Space>
        }
      >
        <Form
          form={form}
          layout="vertical"
          requiredMark={false}
        >
          <Row gutter={16}>
            <Col span={24}>
              <Form.Item
                name="name"  
                label="Name"
                rules={[{ required: true, message: 'Please enter group name' }]}
              >
                <Input placeholder="Please enter group name" />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Drawer>
    </>
  );
};

export default CreateGroupDrawer;
