'use client'

import { Input, FormLabel, Box} from "@chakra-ui/react";

interface Props {
    label?: string
    name?: string,
    value?: number,
    onChange?: (event:React.FormEvent<any>) => void
}

const Number = ( { label, value, name, onChange } : Props) => {
    return (
        <Box maxW='md' minW='3xs'>
            {label && <FormLabel>{label}</FormLabel>}
            <Input 
                type='number'
                value={value}
                name={name}
                onChange={onChange}
            />
        </Box>
    );
};

export default Number;