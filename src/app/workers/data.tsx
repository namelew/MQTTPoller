'use client'
import { Tr, Td, Checkbox } from "@chakra-ui/react";
import { IWorker } from "interfaces/IWorker";

interface Props {
    workers?:IWorker[],
    selectedWorkers:string[],
    selectWorkers: React.Dispatch<React.SetStateAction<string[]>>
}

const Workers = ( { workers, selectedWorkers, selectWorkers } : Props) => {
    const handleSelectRow = (id:string) => {
        if (selectedWorkers.includes(id)) {
            selectWorkers(selectedWorkers.filter((rowId) => rowId !== id));
        } else {
            selectWorkers([...selectedWorkers, id]);
        }
    };

    return (
        <>
            {workers?.map((worker) => (
            <Tr key={worker.id}>
                <Td>
                <Checkbox
                    isChecked={selectedWorkers.includes(worker.id)}
                    onChange={() => handleSelectRow(worker.id)}
                />
                </Td>
                <Td>{worker.id}</Td>
                <Td>{worker.online ? 'Online' : 'Offline'}</Td>
            </Tr>
            ))}
        </>
    );
}

export default Workers;