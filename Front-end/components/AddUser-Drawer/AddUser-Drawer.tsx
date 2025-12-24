import React, { useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Col, Drawer, Form, Input, Row, Space } from 'antd';

const AddUserDrawer: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();

  const showDrawer = () => {
    setOpen(true);
  };

  const onClose = () => {
    setOpen(false);
  };

  const saveUser = async () => {
    try {
      const values = await form.validateFields();
      console.log("Email:", values.Email);

      onClose();
      form.resetFields();
    } catch (error) {
      console.log("Validation error:", error);
    }
  };

  return (
    <>
      <Button type="primary" className='addMemberBtn' onClick={showDrawer} icon={<PlusOutlined />}>
        Add Member
      </Button>

      <Drawer
        title="Add User"
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
            <Button className='addMemberBtn' onClick={saveUser} type="primary">
              Invite User
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
                name="Email"
                label="Email"
                rules={[{ required: true, message: 'Please enter user Email' }]}
              >
                <Input placeholder="Please enter user Email" />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Drawer>
    </>
  );
};

export default AddUserDrawer;
