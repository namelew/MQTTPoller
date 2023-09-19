'use client'
import { Button, Flex, Table, Tr, Th, Thead, Tbody, Checkbox} from "@chakra-ui/react";
import { useState } from "react";
import { IWorker } from "../interfaces/IWorker";
import Workers from "./data";

interface Props {
    workers?:IWorker[],
}

const WorkersTable = ( { workers } : Props) => {
    const [selectAll, setSelectAll] = useState(false);
    const [selectedWorkers, setSelectedWorkers] = useState<string[]>([]);

    const handleSelectAll = () => {
        if (workers) {
            if (selectAll) {
                setSelectedWorkers([]);
            } else {
                setSelectedWorkers(workers.map(worker => worker.id));
            }
            setSelectAll(!selectAll);
        }
    };

    const handleSelectRow = (id:string) => {
        if (selectedWorkers.includes(id)) {
            setSelectedWorkers(selectedWorkers.filter((rowId) => rowId !== id));
        } else {
            setSelectedWorkers([...selectedWorkers, id]);
        }
    };

    const startExperiment = () => {

    };

    return (
        <>
            <Flex justifyContent={'flex-end'}>
                <Button onClick={() => startExperiment()}>Iniciar Experimento</Button>
            </Flex>
            <Table variant='simple'>
                <Thead>
                    <Tr>
                    <Th>
                        <Checkbox isChecked={selectAll} onChange={handleSelectAll} />
                    </Th>
                    <Th>ID</Th>
                    <Th>Status</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    <Workers 
                        workers={workers}
                        selectWorkers={setSelectedWorkers}
                        selectedWorkers={selectedWorkers}
                    />
                </Tbody>
            </Table>
        </>
    );
};

export default WorkersTable;