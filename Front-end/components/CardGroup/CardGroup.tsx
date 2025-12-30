import './style.css'
import React from 'react';
import axiosInstance from '../../axios';
import { useGroups } from '../../context/GroupsContext';



interface CardGroupProps  {
  Id: number;
  Name: string;
} 

const CardGroup: React.FC<CardGroupProps> = ({Id, Name }) => {
    const { refreshGroups } = useGroups();

    const DeleteGroup = async (Id: number) =>{
      try{
        const response = await axiosInstance.delete(
          `/api/group/deleteGroup/${Id}`
        );
        refreshGroups(); 
        console.log(response)
      }catch(error){
        console.log(error)
      }
    }
    
    return (
    <>
      <div id={Id.toString()}  className='ContainerGroupCard'>
        <div className='container-data'>
          <h2>{Name}</h2>
          <div className='Buttons-CardGroup'>
            <button className='open-button'>Open</button>
            <button className='delete-button' onClick={() => DeleteGroup(Id)}>Delete</button>
          </div>
        </div>
      </div>
    </>
    );
};

export default CardGroup;
