'use client'

import { Box } from "@chakra-ui/react";

const AppDefault = ( { children } : { children: React.ReactNode } ) => {
    return (
        <Box w='100%' height='100%' p={10}>
            { children }
        </Box>
    );
};

export default AppDefault;