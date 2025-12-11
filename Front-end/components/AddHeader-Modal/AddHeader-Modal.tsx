import React, { useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Col, Drawer, Form, Input, Row, Space } from 'antd';
import './style.css'

const AddHeaderModal: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();

  const showDrawer = () => {
    setOpen(true);
  };

  const onClose = () => {
    setOpen(false);
  };

  const saveEndpoint = async () => {
    try {
      const values = await form.validateFields();
      console.log(values.headers);

      onClose();
      form.resetFields();
    } catch (error) {
      console.log("Validation error:", error);
    }
  };

  return (
    <>
      <Button type="primary" className='addMemberBtn' onClick={showDrawer} icon={<PlusOutlined />}>
        Add Header
      </Button>

      <Drawer
        title="Add Header"
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
            <Button className='addMemberBtn' onClick={saveEndpoint} type="primary">
              Add Header
            </Button>
          </Space>
        }
      >
       <Form form={form} layout="vertical" requiredMark={false}>
          <Form.List name="headers">
            {(fields, { add, remove }) => (
              <>
                {fields.map(({ key, name, ...restField }) => (
                  <Row gutter={16} key={key} style={{ marginBottom: 10 }}>
                    <Col span={11}>
                      <Form.Item
                        {...restField}
                        name={[name, "name"]}
                        label="Name"
                        rules={[{ required: true, message: "Please enter name" }]}
                      >
                        <Input placeholder="Header name" />
                      </Form.Item>
                    </Col>

                    <Col span={11}>
                      <Form.Item
                        {...restField}
                        name={[name, "value"]}
                        label="Value"
                        rules={[{ required: true, message: "Please enter value" }]}
                      >
                        <Input placeholder="Header value" />
                      </Form.Item>
                    </Col>

                    <Col span={2} style={{ display: "flex", alignItems: "center" }}>
                      <Button danger onClick={() => remove(name)}>
                        X
                      </Button>
                    </Col>
                  </Row>
                ))}

                <Form.Item>
                  <Button
                    type="dashed"
                    onClick={() => add()}
                    block
                    icon={<PlusOutlined />}
                  >
                    Add another header
                  </Button>
                </Form.Item>
              </>
            )}
          </Form.List>
        </Form>

      </Drawer>
    </>
  );
};

export default AddHeaderModal;
