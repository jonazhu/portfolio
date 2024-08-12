# %% [markdown]
# # Final Project
# ## Jonathan Zhu - 03708
# In this project, I will try to write a script to analyze a bacterial genome. The script will identify open reading frames and translate them into proteins, predict their molecular weights and BLAST the five shortest proteins. It will also attempt to write the results to files.

# %%
#import packages
#this cell will only have the import statements for ease later
#i am not even sure if all these will be necessary
import re
import regex
import pyfaidx
import sys
from Bio.Seq import Seq
from Bio.SeqUtils import molecular_weight
from Bio.Blast import NCBIWWW
from Bio import Entrez

# %%
#1: GETTING READING FRAMES
#part 1a: getting the files
#we are going to use a test e.coli file for testing purposes
#but since this is a script, we need to make sure everything is in functions and can be run from command line

#read_file takes in a fasta formate file and returns the sequences corresponding to it.
def read_file(path):
    genes = pyfaidx.Fasta(path)
    records = list(genes.keys())
    seqs = []
    
    for item in records:
        seqs.append(genes[item][0:].seq)

    return seqs


# %%
#part 1b: find the reading frames
#we are going to do a function that will find all potential ORFs, starting with a start TAC and ending with a stop codon ATT, ATC, ACT.
#UPDATE: this function is no longer used for the finding of ORFs as it breaks kernels when used on large sequences. Use find_orfs() instead.
def get_potential_orfs(seq):
    dna = Seq(seq)
    #this uses regular expressions
    #the regex finds a start codon, then any number of characters, and then a stop codon.
    #note that this assumes fully defined bases, so no N bases.
    #also, we use the "regex" module, which supports overlapping
    #we also reverse the strands, given to this in 5' to 3' direction.
    orfs = []
    pos = []
    m = regex.finditer(r"TAC\w*ATT|TAC\w*ATC|TAC\w*ACT", seq[::-1], overlapped=False)
    if m!= None:
        for match in m:
            orfs.append(match.group())
            pos.append(match.start())
    else:
        print("No ORFs found in forward sequence")

    #now check reverse sequence
    m = regex.finditer(r"TAC\w*ATT|TAC\w*ATC|TAC\w*ACT", str(dna.complement()), overlapped=False)
    if m!= None:
        for match in m:
            orfs.append(match.group())
            pos.append(match.start())
            return orfs, pos
    else:
        print("No ORFs found in reverse sequence")
        return -1 #indicates error

#like above, this function is no longer used to find ORFs but is included so as to show my initial thought process.
def get_orfs(seq):
    potential_orfs, potential_pos = get_potential_orfs(seq)
    true_orfs = []
    true_pos = []
    for o in potential_orfs:
        i = potential_orfs.index(o)

        orf_found = False
        num = 0
        current_orf = ""
        while orf_found == False:
            if num + 3 > len(o):
                orf_found = True
                break
            current_codon = o[num:num+3]
            current_orf += current_codon
            if current_codon == "ATT" or current_codon == "ATC" or current_codon == "ACT":
                true_orfs.append(current_orf)
                true_pos.append(potential_pos[i])
                orf_found = True
            num += 1

    return true_orfs, true_pos

#Part 2: now we need to translate
def translate_seqs(seqs):
    prots = []
    for s in seqs:
        current_seq = Seq(s)
        coding_dna = current_seq.complement()
        mrna = coding_dna.transcribe()
        current_prot = str(mrna.translate(to_stop = True))

        #because they're not useful, I am going to remove the ones that have length 1 (just methionine)
        if len(current_prot) > 1:
            prots.append(current_prot)
    return prots

# %%
#Part 3: Molecular Masses
def get_mol_masses(prots):
    weights = []
    for p in prots:
        current_weight = molecular_weight(p, "protein")
        weights.append(current_weight)

    return weights

# %%
#Part 4: Blasting
def blast_prots(prots, mail = "jazhu@andrew.cmu.edu"):
    Entrez.email = mail

    #error case: if there are no proteins
    if len(prots) < 1:
        print("Warning: No proteins given to blast_prots()")
        return -1 #indicates an error

    #first set up the output file
    save_file = open("prot_blast.xml", "w")
    for p in prots:
        result_handle = NCBIWWW.qblast("blastp", "pdbaa", p)
        save_file.write(result_handle.read())
        save_file.write("\n")
    
    save_file.close()
    result_handle.close()

    return 0 #indicates everything is good

#a helper function to get the smallest five prots
def get_best_prots(prots, n = 5):
    #first, we need all the ones that are 1000 AA or longer
    best_prots = []
    for p in prots:
        if len(p) >= 1000:
            best_prots.append(p)

    #second, we will check if it is less than n
    #but if it is zero, we will continue
    if len(best_prots) <= n and len(best_prots) > 0:
        return best_prots
    
    #case 2: get the smallest n, remove and repeat
    smallest_prots = []
    for i in range(n):
        if len(best_prots) == 0:
            break
        current_min = find_smallest_prot(best_prots)
        smallest_prots.append(current_min)
        best_prots.remove(current_min)
    
    if len(smallest_prots) > 0:
        return smallest_prots
    else:
        #case 3: get the n longest prots
        largest_prots = []
        for i in range(n):
            if len(prots) == 0:
                break
            current_max = find_largest_prot(prots)
            largest_prots.append(current_max)
            prots.remove(current_max)
        return largest_prots

#another helper function
def find_smallest_prot(prots):
    #error handling
    if len(prots) < 1:
        return -1 #indicates an error
    
    current_smallest = prots[0]
    for p in prots:
        if len(p) < len(current_smallest):
            current_smallest = p

    return current_smallest

def find_largest_prot(prots):
    #error handling
    if len(prots) < 1:
        return -1 #indicates an error
    
    current_largest = prots[0]
    for p in prots:
        if len(p) > len(current_largest):
            current_largest = p

    return current_largest


# %%
#Extras: Writing to files
#the file ORFs.fasta will have the orfs found, while the file prots.txt and mw.txt will have the protein sequences and the weights.
def write_orfs(orfs, pos, genomelabel):
    with open("orfs.fasta", "w") as fw:
        current_orf = 0
        for o in orfs:
            fw.write(">ORF " + str(current_orf) + " at position " + str(pos[current_orf]) + " of file " + genomelabel + "\n")
            fw.write(o)
            fw.write("\n")
            current_orf += 1

def write_prots(prots, label):
    with open("prots.txt", "w") as fw:
        fw.write("Proteins read from genome file " + label + ":\n")
        for p in prots:
            fw.write(p)
            fw.write("\n")

def write_mw(mws, label):
    with open("mw.txt", "w") as fw:
        fw.write("Molecular masses of proteins read from genome file " + label + ":\n")
        for m in mws:
            fw.write(str(m))
            fw.write("\n")

#alternative function to read ORFs
def find_orfs(dna):
    #takes a 5' to 3' DNA strand
    #since DNA polymerase reads 3' to 5', we need to flip it AND take the complement to get the forward and reverse strands
    forward_strand = dna[::-1]
    reverse_strand = str(Seq(dna).complement())
    n = len(dna)
    
    #step 1: find start codons
    forward_starts = find_start_codons(forward_strand)
    reverse_starts = find_start_codons(reverse_strand)
    

    #step 2: for every start codon, find the first stop codon, and append this ORF to a main list
    orfs = []
    for s in forward_starts:
        current_orf = ""
        current_num = 0
        start_index = s
        end_index = s+3
        stop_reached = False

        while stop_reached == False:
            current_start = start_index+(current_num * 3)
            current_end = end_index+(current_num * 3)
            current_codon = forward_strand[current_start:current_end]
            current_orf += current_codon
            if current_codon == "ATT" or current_codon == "ATC" or current_codon == "ACT":
                stop_reached = True
                orfs.append(current_orf)
            if current_end >= n:
                stop_reached = True
            current_num += 1

    #now do the same for the other strand
    for s in reverse_starts:
        current_orf = ""
        current_num = 0
        start_index = s
        end_index = s+3
        stop_reached = False

        while stop_reached == False:
            current_start = start_index+(current_num * 3)
            current_end = end_index+(current_num * 3)
            current_codon = reverse_strand[current_start:current_end]
            current_orf += current_codon
            if current_codon == "ATT" or current_codon == "ATC" or current_codon == "ACT":
                stop_reached = True
                orfs.append(current_orf)
            if current_end >= n:
                stop_reached = True
            current_num += 1

    forward_starts.extend(reverse_starts)
    return orfs, forward_starts

#helper functions
#find start codons for a given strand
def find_start_codons(dna):
    #we need TAC codons
    #this one can use regex
    m = re.finditer(r"TAC", dna)
    match_pos = []
    if m != None:
        for i in m:
            match_pos.append(i.start())

    return match_pos

#the main portion of writing the workflow - aka putting all the functions together. 
if __name__ == "__main__":
    file = sys.argv[1] #first argument: the file path. This is mandatory.
    n = len(sys.argv)
    #secondary arguments: additional things.
    mode = "all"
    email = "jazhu@andrew.cmu.edu"
    if n > 2:
        mode = sys.argv[2] #second argument: the mode for which to do. If just want reading frames, set to "rf". If just want prots, set to "prot". If just want molecular weights, set to "mw". If just want blast, set to "blast". If want all, set to "all", or leave blank. If you want all except blast (most sane choice), pick "noblast".
    if n > 3:
        email = sys.argv[3] #third argument: email address for blasting. it can use mine if you don't want to put it in, I get it.

    #the main workflow
    genome_seqs = read_file(file)
    orfs = []
    pos = []
    for s in genome_seqs:
        current_orfs, current_pos = find_orfs(s)
        orfs.extend(current_orfs)
        pos.extend(current_pos)
    prots = translate_seqs(orfs)
    mol_masses = get_mol_masses(prots)
    sample = get_best_prots(prots)

    if mode == "rf":
        write_orfs(orfs, pos, file)
    if mode == "prot":
        write_prots(prots, file)
    if mode == "mw":
        write_mw(mol_masses, file)
    if mode == "blast":
        blast_prots(sample, email)
    if mode == "noblast":
        write_orfs(orfs, pos, file)
        write_prots(prots, file)
        write_mw(mol_masses, file)
    else: #the all case
        write_orfs(orfs, pos, file)
        write_prots(prots, file)
        write_mw(mol_masses, file)
        blast_prots(sample, email)
